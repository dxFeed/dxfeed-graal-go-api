package live

import (
	"path/filepath"
	"slices"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/ipf"
)

const testProfiles = "../testdata/profiles.ipf"

type updates struct {
	ch chan []*events.InstrumentProfile
}

func newUpdates() *updates {
	return &updates{ch: make(chan []*events.InstrumentProfile, 100)}
}

func (u *updates) InstrumentProfilesUpdated(profiles []*events.InstrumentProfile) {
	u.ch <- profiles
}

// next returns the symbols and types of the next update.
func (u *updates) next(t *testing.T) map[string]string {
	t.Helper()
	select {
	case profiles := <-u.ch:
		result := make(map[string]string, len(profiles))
		for _, p := range profiles {
			result[*p.Symbol()] = *p.InstrumentType()
		}
		return result
	case <-time.After(10 * time.Second):
		t.Fatal("no update in 10 seconds")
		return nil
	}
}

func profile(t *testing.T, symbol, instrumentType string) *events.InstrumentProfile {
	t.Helper()
	p := events.NewInstrumentProfile()
	if err := ipf.SetField(p, "SYMBOL", symbol); err != nil {
		t.Fatal(err)
	}
	if err := ipf.SetField(p, "TYPE", instrumentType); err != nil {
		t.Fatal(err)
	}
	return p
}

func symbols(profiles []*events.InstrumentProfile) []string {
	result := make([]string, 0, len(profiles))
	for _, p := range profiles {
		result = append(result, *p.Symbol())
	}
	sort.Strings(result)
	return result
}

func TestCollectorUpdates(t *testing.T) {
	collector, err := NewInstrumentProfileCollector()
	if err != nil {
		t.Fatal(err)
	}
	defer collector.Close()

	first := newUpdates()
	if err := collector.AddUpdateListener(first); err != nil {
		t.Fatal(err)
	}
	if err := collector.UpdateInstrumentProfiles([]*events.InstrumentProfile{profile(t, "AAA", ipf.TypeStock),
		profile(t, "BBB", ipf.TypeStock)}); err != nil {
		t.Fatal(err)
	}
	if got := first.next(t); len(got) != 2 || got["AAA"] != ipf.TypeStock || got["BBB"] != ipf.TypeStock {
		t.Fatalf("first update = %v", got)
	}
	view, err := collector.View()
	if err != nil || !slices.Equal(symbols(view), []string{"AAA", "BBB"}) {
		t.Fatalf("view = %v, %v", symbols(view), err)
	}
	if updated, err := collector.LastUpdateTime(); err != nil || time.Since(updated) > time.Minute {
		t.Errorf("last update time = %v, %v", updated, err)
	}

	// A listener added later first receives all profiles.
	second := newUpdates()
	if err := collector.AddUpdateListener(second); err != nil {
		t.Fatal(err)
	}
	if got := second.next(t); len(got) != 2 {
		t.Fatalf("snapshot of the second listener = %v", got)
	}

	if err := collector.UpdateInstrumentProfile(profile(t, "AAA", ipf.TypeRemoved)); err != nil {
		t.Fatal(err)
	}
	for _, u := range []*updates{first, second} {
		if got := u.next(t); len(got) != 1 || got["AAA"] != ipf.TypeRemoved {
			t.Fatalf("removal = %v", got)
		}
	}
	view, err = collector.View()
	if err != nil || !slices.Equal(symbols(view), []string{"BBB"}) {
		t.Fatalf("view after removal = %v, %v", symbols(view), err)
	}

	if err := collector.RemoveUpdateListener(first); err != nil {
		t.Fatal(err)
	}
	if err := collector.UpdateInstrumentProfile(profile(t, "CCC", ipf.TypeStock)); err != nil {
		t.Fatal(err)
	}
	if got := second.next(t); got["CCC"] != ipf.TypeStock {
		t.Fatalf("update after a removal of a listener = %v", got)
	}
	select {
	case got := <-first.ch:
		t.Errorf("a removed listener got %d profiles", len(got))
	case <-time.After(200 * time.Millisecond):
	}

	if err := collector.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := collector.View(); err == nil {
		t.Error("View after Close did not fail")
	}
	if err := collector.Close(); err != nil {
		t.Errorf("second Close: %v", err)
	}
}

type states struct {
	mu   sync.Mutex
	seen []State
}

func (s *states) StateChanged(_, newState State) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seen = append(s.seen, newState)
}

func (s *states) get() []State {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.seen)
}

func TestConnectionToFile(t *testing.T) {
	collector, err := NewInstrumentProfileCollector()
	if err != nil {
		t.Fatal(err)
	}
	defer collector.Close()
	path, err := filepath.Abs(testProfiles)
	if err != nil {
		t.Fatal(err)
	}
	connection, err := CreateConnection(path+"[update=5s]", collector)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()

	if address, err := connection.Address(); err != nil || address != path {
		t.Errorf("address = %q, %v, want %q", address, err, path)
	}
	if period, err := connection.UpdatePeriod(); err != nil || period != 5*time.Second {
		t.Errorf("update period = %v, %v", period, err)
	}
	if state, err := connection.State(); err != nil || state != NotConnected {
		t.Errorf("state before Start = %v, %v", state, err)
	}
	listener := &states{}
	if err := connection.AddStateChangeListener(listener); err != nil {
		t.Fatal(err)
	}
	updates := newUpdates()
	if err := collector.AddUpdateListener(updates); err != nil {
		t.Fatal(err)
	}

	if err := connection.Start(); err != nil {
		t.Fatal(err)
	}
	completed, err := connection.WaitUntilCompleted(10 * time.Second)
	if err != nil || !completed {
		t.Fatalf("WaitUntilCompleted = %v, %v", completed, err)
	}
	if got := updates.next(t); len(got) != 9 || got["TEST"] != ipf.TypeStock || got[".TEST261218C10"] != ipf.TypeOption {
		t.Errorf("update from the file = %v", got)
	}
	view, err := collector.View()
	if err != nil || len(view) != 9 {
		t.Errorf("view = %d profiles, %v", len(view), err)
	}
	if state, err := connection.State(); err != nil || state != Completed {
		t.Errorf("state = %v, %v", state, err)
	}
	if modified, err := connection.LastModified(); err != nil || modified.IsZero() {
		t.Errorf("last modified = %v, %v", modified, err)
	}
	if err := connection.Start(); err == nil {
		t.Error("a second Start did not fail")
	}

	if err := connection.Close(); err != nil {
		t.Fatal(err)
	}
	if state, err := connection.State(); err != nil || state != Closed {
		t.Errorf("state after Close = %v, %v", state, err)
	}
	if got, want := listener.get(), []State{Connecting, Connected, Completed, Closed}; !slices.Equal(got, want) {
		t.Errorf("states = %v, want %v", got, want)
	}
	if _, err := connection.Address(); err == nil {
		t.Error("Address after Close did not fail")
	}
}

func TestConnectionErrors(t *testing.T) {
	if _, err := CreateConnection("x.ipf", nil); err == nil {
		t.Error("a nil collector did not fail")
	}
	collector, err := NewInstrumentProfileCollector()
	if err != nil {
		t.Fatal(err)
	}
	connection, err := CreateConnection("x.ipf", collector)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()
	if err := connection.SetUpdatePeriod(-time.Second); err == nil {
		t.Error("a negative update period did not fail")
	}
	if err := connection.SetUpdatePeriod(1500 * time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if period, _ := connection.UpdatePeriod(); period != 1500*time.Millisecond {
		t.Errorf("update period = %v", period)
	}
	if State(9).String() != "State(9)" || Completed.String() != "COMPLETED" {
		t.Error("State.String is wrong")
	}
	if err := collector.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateConnection("x.ipf", collector); err == nil {
		t.Error("a closed collector did not fail")
	}
}
