package native

/*
#include "dxfg_api.h"
extern void OnSymbolsAdded(graal_isolatethread_t *thread, dxfg_symbol_list *symbols, void *user_data);
extern void OnSymbolsRemoved(graal_isolatethread_t *thread, dxfg_symbol_list *symbols, void *user_data);
extern void OnSubscriptionClosed(graal_isolatethread_t *thread, void *user_data);
extern void OnConfigurationChanged(graal_isolatethread_t *thread, void *user_data);
*/
import "C"
import (
	"errors"
	"log"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
)

//export OnSymbolsAdded
func OnSymbolsAdded(thread *C.graal_isolatethread_t, symbols *C.dxfg_symbol_list, userData unsafe.Pointer) {
	defer recoverCallback("subscription change listener")
	if listener, ok := Restore(userData).(osub.ObservableSubscriptionChangeListener); ok {
		listener.SymbolsAdded(changedSymbols(symbols))
	}
}

//export OnSymbolsRemoved
func OnSymbolsRemoved(thread *C.graal_isolatethread_t, symbols *C.dxfg_symbol_list, userData unsafe.Pointer) {
	defer recoverCallback("subscription change listener")
	if listener, ok := Restore(userData).(osub.ObservableSubscriptionChangeListener); ok {
		listener.SymbolsRemoved(changedSymbols(symbols))
	}
}

//export OnSubscriptionClosed
func OnSubscriptionClosed(thread *C.graal_isolatethread_t, userData unsafe.Pointer) {
	defer recoverCallback("subscription change listener")
	if listener, ok := Restore(userData).(osub.ObservableSubscriptionChangeListener); ok {
		listener.SubscriptionClosed()
	}
}

//export OnConfigurationChanged
func OnConfigurationChanged(thread *C.graal_isolatethread_t, userData unsafe.Pointer) {
	defer recoverCallback("subscription change listener")
	if listener, ok := Restore(userData).(osub.ConfigurationChangeListener); ok {
		listener.ConfigurationChanged()
	}
}

// changedSymbols copies the symbols of a change notification; the list belongs to the SDK. A symbol that does not
// convert is logged and left out, so that the listener still learns of the others.
func changedSymbols(list *C.dxfg_symbol_list) []osub.Symbol {
	if list == nil || list.size <= 0 || list.elements == nil {
		return []osub.Symbol{}
	}
	elements := unsafe.Slice((**C.dxfg_v4_symbol_t)(unsafe.Pointer(list.elements)), int(list.size))
	symbols := make([]osub.Symbol, 0, len(elements))
	for _, s := range elements {
		symbol, err := goSymbol(s)
		if err != nil {
			log.Printf("dxfeed: subscription change listener: a symbol is skipped: %v", err)
			continue
		}
		symbols = append(symbols, symbol)
	}
	return symbols
}

// nativeListener is a native listener with the Go value its callbacks restore from the user data.
type nativeListener struct {
	handle   Handler
	userData unsafe.Pointer
}

// Release releases the native listener: the SDK does not call it after Release returns, except on the goroutine that
// releases it from the listener itself, so the Go value is forgotten too. Later calls do nothing.
func (l *nativeListener) Release() error {
	if l.handle.Ptr() == nil {
		return nil
	}
	if err := l.handle.Free(); err != nil {
		// The SDK may still call the listener: keep the Go value for it.
		return err
	}
	Unref(l.userData)
	return nil
}

// ChangeListener is a native listener of the changes of a subscription, created with NewChangeListener.
type ChangeListener struct {
	nativeListener
}

// NewChangeListener creates a native listener that passes the changes to listener, and the configuration changes if
// listener is also an osub.ConfigurationChangeListener.
func NewChangeListener(listener osub.ObservableSubscriptionChangeListener) (*ChangeListener, error) {
	userData := Save(listener)
	var configurationChanged *[0]byte
	if _, ok := listener.(osub.ConfigurationChangeListener); ok {
		configurationChanged = (*[0]byte)(C.OnConfigurationChanged)
	}
	var l *C.dxfg_v4_observable_subscription_change_listener_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_ObservableSubscriptionChangeListener_newWithConfigurationChanged(thread.ptr,
			(*[0]byte)(C.OnSymbolsAdded), (*[0]byte)(C.OnSymbolsRemoved), (*[0]byte)(C.OnSubscriptionClosed),
			configurationChanged, userData, &l))
	})
	if l == nil {
		Unref(userData)
		if err == nil {
			err = errors.New("native: the SDK returned no change listener")
		}
		return nil, err
	}
	cl := &ChangeListener{nativeListener{handle: NewJavaHandle(unsafe.Pointer(l)), userData: userData}}
	if err != nil {
		return nil, errors.Join(err, cl.Release())
	}
	return cl, nil
}

func (cl *ChangeListener) ptr() *C.dxfg_v4_observable_subscription_change_listener_t {
	return (*C.dxfg_v4_observable_subscription_change_listener_t)(cl.handle.Ptr())
}

// AddChangeListener adds the change listener to the subscription.
func (s DXFeedSubscription) AddChangeListener(cl *ChangeListener) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_addChangeListener(thread.ptr, s.ptr, cl.ptr()))
	})
}

// RemoveChangeListener removes the change listener from the subscription; it does not release it.
func (s DXFeedSubscription) RemoveChangeListener(cl *ChangeListener) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_DXFeedSubscription_removeChangeListener(thread.ptr, s.ptr, cl.ptr()))
	})
}
