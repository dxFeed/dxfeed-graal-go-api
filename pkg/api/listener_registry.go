package api

import (
	"log"
	"sync"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
)

// listenerRegistry keeps the listeners added to a subscription with their native listeners of type N. An entry is in
// the list before its native listener is attached, so a listener may remove itself from its first call.
type listenerRegistry[L comparable, N any] struct {
	// kind names the listeners in the log, such as "event listener".
	kind    string
	mu      sync.Mutex
	entries []*registryEntry[L, N]
	closed  bool
}

// changeListenerRegistry keeps the change listeners of a subscription.
type changeListenerRegistry = listenerRegistry[osub.ObservableSubscriptionChangeListener, *native.ChangeListener]

func newChangeListenerRegistry() changeListenerRegistry {
	return changeListenerRegistry{kind: "subscription change listener"}
}

type registryEntry[L comparable, N any] struct {
	listener L
	native   N
	attached bool
	removed  bool
}

// add adds the listener: attach creates its native listener and adds it to the subscription, detach removes and
// releases a native listener. A closed registry ignores the listener.
func (r *listenerRegistry[L, N]) add(listener L, attach func() (N, error), detach func(N) error) error {
	entry := &registryEntry[L, N]{listener: listener}
	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return nil
	}
	r.entries = append(r.entries, entry)
	r.mu.Unlock()

	native, err := attach()
	r.mu.Lock()
	if err != nil {
		r.removeEntry(entry)
		r.mu.Unlock()
		return err
	}
	if !entry.removed {
		entry.native = native
		entry.attached = true
		r.mu.Unlock()
		return nil
	}
	r.mu.Unlock()
	// remove or close came while the native listener was being attached.
	r.logDetach(detach(native))
	return nil
}

// remove removes the first entry of the listener and detaches its native listener. If add of that entry is still
// attaching it, add detaches it when done.
func (r *listenerRegistry[L, N]) remove(listener L, detach func(N) error) {
	r.mu.Lock()
	var native N
	attached := false
	for _, e := range r.entries {
		if e.listener == listener {
			native, attached = e.native, e.attached
			r.removeEntry(e)
			break
		}
	}
	r.mu.Unlock()
	if attached {
		r.logDetach(detach(native))
	}
}

// close marks the registry closed, empties it and returns the attached native listeners.
func (r *listenerRegistry[L, N]) close() []N {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closed = true
	natives := make([]N, 0, len(r.entries))
	for _, e := range r.entries {
		e.removed = true
		if e.attached {
			natives = append(natives, e.native)
		}
	}
	r.entries = nil
	return natives
}

// removeEntry marks the entry removed and takes it out of the list; r.mu must be held.
func (r *listenerRegistry[L, N]) removeEntry(entry *registryEntry[L, N]) {
	entry.removed = true
	for i, e := range r.entries {
		if e == entry {
			r.entries = append(r.entries[:i:i], r.entries[i+1:]...)
			return
		}
	}
}

func (r *listenerRegistry[L, N]) logDetach(err error) {
	if err != nil {
		log.Printf("dxfeed: remove %s: %v", r.kind, err)
	}
}
