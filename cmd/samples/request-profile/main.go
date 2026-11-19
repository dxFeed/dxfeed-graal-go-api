// Command request-profile connects to a running publish-profiles sample, prints the profile of a ":TEST" symbol and
// exits.
//
//	go run ./cmd/samples/request-profile -address localhost:7700 -symbol A:TEST
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
)

func main() {
	address := flag.String("address", "localhost:7700", "the address of the publish-profiles sample")
	symbol := flag.String("symbol", "A:TEST", "the symbol, ending with :TEST")
	timeout := flag.Duration("timeout", 5*time.Second, "how long to wait for the profile")
	flag.Parse()

	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		panic(err)
	}
	defer func() { _ = endpoint.Close() }()
	if err := endpoint.Connect(*address); err != nil {
		panic(err)
	}
	feed, err := endpoint.GetFeed()
	if err != nil {
		panic(err)
	}
	promise, err := feed.GetLastEventPromise(eventcodes.Profile, osub.String(*symbol))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	profile, err := promise.Await(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "no profile of %s: %v\n", *symbol, err)
		os.Exit(1)
	}
	fmt.Println(profile)
}
