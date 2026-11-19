// Command convert-tape-file converts one tape file into another, with optional processing or filtering of the events
// in between. ConvertTapeFile.in in this directory is a sample input.
//
//	cd cmd/samples/convert-tape-file
//	go run . -in "file:ConvertTapeFile.in[readAs=stream_data,speed=max]" \
//	  -out "tape:ConvertTapeFile.out[saveAs=stream_data,format=text]"
package main

import (
	"context"
	"flag"
	"sync"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

// eventListener adapts a function to common.EventListener.
type eventListener func(events []interface{})

func (f eventListener) Update(events []interface{}) { f(events) }

// newTapeEndpoint creates an endpoint of a stream role, which neither conflates nor drops events, that receives every
// event with the wildcard subscription and keeps the event times of the tape.
func newTapeEndpoint(role common.Role) *api.DXEndpoint {
	endpoint, err := api.NewEndpointBuilder().
		WithRole(role).
		WithProperty("dxfeed.wildcard.enable", "true").
		WithProperty("dxendpoint.eventTime", "true").
		Build()
	if err != nil {
		panic(err)
	}
	return endpoint
}

func main() {
	in := flag.String("in", "file:ConvertTapeFile.in[readAs=stream_data,speed=max]", "the input tape address")
	out := flag.String("out", "tape:ConvertTapeFile.out[saveAs=stream_data,format=text]", "the output tape address")
	flag.Parse()

	input := newTapeEndpoint(api.StreamFeed)
	defer func() { _ = input.Close() }()
	output := newTapeEndpoint(api.StreamPublisher)
	defer func() { _ = output.Close() }()

	feed, err := input.GetFeed()
	if err != nil {
		panic(err)
	}
	publisher, err := output.GetPublisher()
	if err != nil {
		panic(err)
	}
	// All the event types of the endpoint; the set of processed types may be limited if needed.
	types, err := input.EventTypes()
	if err != nil {
		panic(err)
	}
	subscription, err := feed.CreateSubscription(types...)
	if err != nil {
		panic(err)
	}
	defer func() { _ = subscription.Close() }()
	// A panic in a listener does not end the program: keep the first error of Publish and report it at the end.
	var publishErr error
	var publishMu sync.Mutex
	err = subscription.AddListener(eventListener(func(events []interface{}) {
		// Here the events may be changed, removed or added: for example, set the event time of each event 1 hour later
		// with event.SetEventTime(event.EventTime() + 3_600_000).
		if err := publisher.Publish(events); err != nil {
			publishMu.Lock()
			if publishErr == nil {
				publishErr = err
			}
			publishMu.Unlock()
		}
	}))
	if err != nil {
		panic(err)
	}
	// All the symbols; the set of processed symbols may be limited if needed.
	if err := subscription.AddSymbol(osub.NewWildcardSymbol()); err != nil {
		panic(err)
	}

	// Start writing the output tape BEFORE reading the input one.
	if err := output.Connect(*out); err != nil {
		panic(err)
	}
	if err := input.Connect(*in); err != nil {
		panic(err)
	}

	ctx := context.Background()
	// Wait until all the data is read and processed, then close the input endpoint.
	if err := input.AwaitNotConnected(ctx); err != nil {
		panic(err)
	}
	if err := input.CloseAndAwaitTermination(ctx); err != nil {
		panic(err)
	}
	// Wait until all the data is processed and written, then close the output endpoint.
	if err := output.AwaitProcessed(ctx); err != nil {
		panic(err)
	}
	if err := output.CloseAndAwaitTermination(ctx); err != nil {
		panic(err)
	}
	publishMu.Lock()
	defer publishMu.Unlock()
	if publishErr != nil {
		panic(publishErr)
	}
}
