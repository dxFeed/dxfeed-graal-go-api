package api

import (
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/profile"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

// TestPublishFreesEvents checks that Publish frees the C memory of the published events.
func TestPublishFreesEvents(t *testing.T) {
	_, publisher := newLocalHub(t)
	before := mappers.LiveAllocations()
	for i := 0; i < 100; i++ {
		p := profile.NewProfile("PUBLISH/MEMORY")
		description := "description"
		p.SetDescription(&description)
		if err := publisher.Publish([]interface{}{quote.NewQuote("PUBLISH/MEMORY"), p}); err != nil {
			t.Fatalf("publish: %v", err)
		}
	}
	if after := mappers.LiveAllocations(); after != before {
		t.Fatalf("%d C allocations left after publishing", after-before)
	}
}
