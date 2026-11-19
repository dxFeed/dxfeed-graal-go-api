// Command write-tape-file publishes two quotes into a tape file in the text format.
//
//	go run ./cmd/samples/write-tape-file -address "tape:WriteTapeFile.out.txt[format=text]"
//
// Remove [format=text] or change it to [format=binary] to write the binary format.
package main

import (
	"context"
	"flag"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

func main() {
	address := flag.String("address", "tape:WriteTapeFile.out.txt[format=text]", "the tape address to write to")
	flag.Parse()

	endpoint, err := api.NewEndpointBuilder().
		// The tape connector receives every event only with the wildcard subscription.
		WithProperty("dxfeed.wildcard.enable", "true").
		WithRole(api.Publisher).
		Build()
	if err != nil {
		panic(err)
	}
	defer func() { _ = endpoint.Close() }()
	if err := endpoint.Connect(*address); err != nil {
		panic(err)
	}
	publisher, err := endpoint.GetPublisher()
	if err != nil {
		panic(err)
	}

	quote1 := quote.NewQuote("TEST1")
	quote1.SetBidPrice(10.1)
	quote1.SetAskPrice(10.2)
	quote2 := quote.NewQuote("TEST2")
	quote2.SetBidPrice(17.1)
	quote2.SetAskPrice(17.2)
	if err := publisher.Publish([]interface{}{quote1, quote2}); err != nil {
		panic(err)
	}

	// Wait until all the data is written, then close the endpoint and wait until it is closed.
	ctx := context.Background()
	if err := endpoint.AwaitProcessed(ctx); err != nil {
		panic(err)
	}
	if err := endpoint.CloseAndAwaitTermination(ctx); err != nil {
		panic(err)
	}
}
