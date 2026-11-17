package api

import (
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

type unsupportedEvent struct{}

func (unsupportedEvent) Type() eventcodes.EventCode { return eventcodes.Series }

func TestPublishUnsupportedEventReturnsError(t *testing.T) {
	endpoint, err := NewEndpoint(Publisher)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = endpoint.Close() }()
	publisher, err := endpoint.GetPublisher()
	if err != nil {
		t.Fatal(err)
	}
	if err := publisher.Publish([]interface{}{unsupportedEvent{}}); err == nil {
		t.Fatal("publishing an event without a native mapper must return an error")
	}
}

func TestSubscribeUnsupportedEventReturnsError(t *testing.T) {
	endpoint, err := NewEndpoint(Feed)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = endpoint.Close() }()
	feed, err := endpoint.GetFeed()
	if err != nil {
		t.Fatal(err)
	}
	sub, err := feed.CreateSubscription(eventcodes.Quote, eventcodes.Series)
	if err == nil {
		t.Fatal("subscribing to an event without a native mapper must return an error")
	}
	if sub != nil {
		t.Fatal("a failed CreateSubscription must return a nil subscription")
	}
}
