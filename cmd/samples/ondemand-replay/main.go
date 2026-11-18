// Command ondemand-replay replays the quotes of Accenture (ACN), which dropped under $1 in the "Flash Crash" of
// May 6, 2010 from 14:47:48 to 14:48:02 EST, from the dxFeed on-demand history data replay service.
package main

import (
	"fmt"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
)

type PrintEvents func(events []interface{})

func (pr PrintEvents) Update(events []any) {
	pr(events)
}

const timeLayout = "2006-01-02 15:04:05.000 MST"

func main() {
	endpoint, err := api.NewEndpointWithProperties(api.OnDemandFeed, map[string]string{
		// The replayed events keep the time when they happened.
		"dxendpoint.eventTime": "true",
		"dxfeed.user":          "demo",
		"dxfeed.password":      "demo",
	})
	if err != nil {
		panic(err)
	}
	defer func() { _ = endpoint.Close() }()

	onDemand, err := endpoint.OnDemandService()
	if err != nil {
		panic(err)
	}
	feed, err := endpoint.GetFeed()
	if err != nil {
		panic(err)
	}
	subscription, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		panic(err)
	}
	defer func() { _ = subscription.Close() }()
	est := time.FixedZone("EST", -5*60*60)
	err = subscription.AddListener(PrintEvents(func(events []interface{}) {
		for _, event := range events {
			if q, ok := event.(*quote.Quote); ok {
				fmt.Printf("%s : %s bid %v / ask %v\n", time.UnixMilli(q.EventTime()).In(est).Format(timeLayout),
					*q.EventSymbol(), q.BidPrice(), q.AskPrice())
			}
		}
	}))
	if err != nil {
		panic(err)
	}
	if err = subscription.AddSymbol(osub.String("ACN")); err != nil {
		panic(err)
	}

	// The address of the on-demand service: the "(ondemand:...)" component makes the replay available.
	if err = endpoint.Connect("(ondemand:demo.dxfeed.com:7680)"); err != nil {
		panic(err)
	}
	from := time.Date(2010, 5, 6, 14, 47, 48, 0, est)
	to := time.Date(2010, 5, 6, 14, 48, 2, 0, est)
	if err = onDemand.Replay(from); err != nil {
		panic(err)
	}
	for {
		now, err := onDemand.Time()
		if err != nil {
			panic(err)
		}
		if !now.Before(to) {
			break
		}
		state, _ := endpoint.State()
		fmt.Printf("Current state is %s, on-demand time is %s\n", state, now.In(est).Format(timeLayout))
		time.Sleep(time.Second)
	}
}
