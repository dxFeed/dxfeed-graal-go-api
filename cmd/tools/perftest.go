package main

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/timeandsale"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/parser"
	"math"
	"os"
	"sync"
	"time"
	"unsafe"
)

type PerfTest struct{}

var hash uintptr

func (c PerfTest) ShortDescription() string {
	return "Connects to specified address and calculates performance counters."
}

func (c PerfTest) Run(args []string) {
	dxarguments := DXArguments{args}
	arguments := dxarguments.arguments()
	if len(arguments) < 3 {
		fmt.Println(`
		Connects to the specified address(es) and calculates performance counters (events per second, cpu usage, etc).
		
		Usage:
		path_to_app <address> <types> <symbols> [<options>]
		
		Where:
		
		address (pos. 0)  Required. The address(es) to connect to retrieve data (see "Help address").
					For Token-Based Authorization, use the following format: "<address>:<port>[login=entitle:<token>]".
		types (pos. 1)    Required. Comma-separated list of dxfeed event types (e.g. Quote, TimeAndSale).
		symbols (pos. 2)  Required. Comma-separated list of symbol names to get events for (e.g. "IBM, AAPL, MSFT").
		--force-stream    Enforces a streaming contract for subscription. The StreamFeed role is used instead of Feed.`)
		os.Exit(0)
	}
	address := arguments[0]
	types := parser.ParseEventTypes(arguments[1])
	symbols := parser.ParseSymbols(arguments[2])

	err := perf(address, types, symbols, dxarguments.forceStream())
	if err != nil {
		fmt.Printf("Error during dump: %v", err)
	}
}

func perf(address string, types []eventcodes.EventCode, symbols []any, forceStream bool) error {
	role := api.Feed
	if forceStream {
		role = api.StreamFeed
	}
	endpoint, err := api.CreateEndpoint(role)
	if err != nil {
		return fmt.Errorf("CreateEndpoint: %we", err)
	}
	err = endpoint.Connect(address)
	feed, err := endpoint.GetFeed()
	if err != nil {
		return fmt.Errorf("GetFeed: %we", err)
	}
	sub, err := feed.CreateSubscription(types...)
	if err != nil {
		return fmt.Errorf("CreateSubscription: %we", err)
	}
	d := &diag{}

	err = sub.AddListener(PrintEvents(func(eventsList []interface{}) {
		d.mu.Lock()
		d.addListenerCounter(1)
		d.addEventCounter(len(eventsList))
		for _, event := range eventsList {
			switch v := event.(type) {
			case *quote.Quote:
				fmt.Printf("%s\n", v.String())
			case *timeandsale.TimeAndSale:
				hash += uintptr(unsafe.Pointer(&v))
			default:
				fmt.Printf("Unsupported event %T!\n", v)
			}
		}
		d.mu.Unlock()
	}))
	if err != nil {
		return fmt.Errorf("AddListener: %we", err)
	}
	for _, symbol := range symbols {
		err = sub.AddSymbol(symbol)
		if err != nil {
			return fmt.Errorf("AddSymbol: %we", err)

		}
	}

	go func() {
		for {
			time.Sleep(2 * time.Second)
			d.PrintDiag()
		}
	}()

	time.Sleep(time.Duration(math.MaxInt64))
	fmt.Println(hash)
	return nil
}

type PrintEvents func(events []interface{})

func (pr PrintEvents) Update(events []any) {
	pr(events)
}

type diag struct {
	listenerCounter int
	eventCounter    int
	mu              sync.Mutex
}

func (d *diag) addListenerCounter(i int) {
	d.listenerCounter += i
}

func (d *diag) addEventCounter(i int) {
	d.eventCounter += i
}

func (d *diag) PrintDiag() {
	d.mu.Lock()
	eventPerSec := float64(d.eventCounter) / 2.0
	listenerCallsPerSec := float64(d.listenerCounter) / 2.0
	fmt.Println("----------------------------------------------")
	fmt.Printf("Rate of events (avg)           : %f (events/s)\n", eventPerSec)
	fmt.Printf("Rate of listener calls         : %f (calls/s)\n", listenerCallsPerSec)
	fmt.Printf("Number of events in call (avg) : %f (events)\n", eventPerSec/listenerCallsPerSec)
	d.eventCounter = 0
	d.listenerCounter = 0
	d.mu.Unlock()
}
