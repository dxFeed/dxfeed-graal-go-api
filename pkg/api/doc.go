// Package api connects to dxFeed: an endpoint ([DXEndpoint]) gives a feed ([DXFeed]) to subscribe to market events
// and a publisher ([DXPublisher]) to publish them. It wraps the dxFeed Graal Native SDK, the Java API of dxFeed
// compiled to a native library, and follows the dxFeed Java and C++ APIs.
//
// # Installation
//
// The module links the native library of the SDK with cgo, so a program needs a C compiler and the library: install
// it with the install-sdk command of the module and set CGO_LDFLAGS and the library path it prints. The guide
// docs/quick-start.md of the module describes the steps for Windows, Linux and macOS.
//
// # A minimal program
//
//	endpoint, err := api.NewEndpoint(api.Feed)
//	if err != nil {
//		return err
//	}
//	defer endpoint.Close()
//	if err := endpoint.Connect("demo.dxfeed.com:7300"); err != nil {
//		return err
//	}
//	feed, err := endpoint.GetFeed()
//	if err != nil {
//		return err
//	}
//	subscription, err := feed.CreateSubscription(eventcodes.Quote)
//	if err != nil {
//		return err
//	}
//	defer subscription.Close()
//	err = subscription.AddListener(listener) // a common.EventListener; its Update gets *quote.Quote values
//	if err != nil {
//		return err
//	}
//	if err := subscription.AddSymbols(osub.Strings("AAPL", "IBM")...); err != nil {
//		return err
//	}
//	// The events come to the listener while the program runs: wait here (a signal, a timer) before returning.
//
// # Endpoints and roles
//
// An endpoint is a connection to dxFeed or to a local hub. Its role ([common.Role]) is fixed at creation:
//
//   - [Feed] receives events and conflates fast updates of the same symbol; [StreamFeed] receives every event;
//     [OnDemandFeed] replays history.
//   - [Publisher] and [StreamPublisher] send events to the address they connect to.
//   - [LocalHub] connects a feed and a publisher in one process, without a network.
//
// [NewEndpoint] creates an endpoint with the default properties; [NewEndpointBuilder] sets properties, a properties
// file, a name and the executor of the notifications. [DXEndpoint.Connect] takes an address such as
// "demo.dxfeed.com:7300" or a file; an authorization token goes into the address as
// "host:port[login=entitle:<token>]". [DXEndpoint.AddListener] reports the connection
// state, [DXEndpoint.Connectors] and [DXEndpoint.LastConnectionError] the connections and their errors.
//
// # Subscriptions and symbols
//
// [DXFeed.CreateSubscription] creates a [DXFeedSubscription] of one or more event types (package eventcodes). Its
// listeners get the events of the subscribed symbols in batches. The methods that take symbols take [osub.Symbol],
// a sealed interface: a plain symbol is [osub.String] ([osub.Strings] converts several), and the wildcard, time
// series and indexed event symbols are types of package osub too. A string or any other type is rejected at compile
// time; package osub explains why.
//
// # Event types
//
// The events are pointers to the types of the packages under events, such as *quote.Quote, *trade.Trade or
// *order.Order. Package events describes the kinds of events: lasting events (the last value of a symbol), indexed
// events (several values of a symbol, such as the orders of a book) and time series events (history by time, such as
// candles).
//
// # Time series
//
// A subscription to a time series event with an [osub.TimeSeriesSubscriptionSymbol] gets the events of the symbol
// from the given time on, then the new ones as they come. For candles the symbol is a string with the period and the
// other attributes of the aggregation in braces, such as "AAPL{=1d}" (package candle).
//
// # Publishing
//
// [DXEndpoint.GetPublisher] of a publisher or local hub endpoint gives a [DXPublisher]; [DXPublisher.Publish] sends
// events to the subscribers. [DXEndpoint.AwaitProcessed] waits until the published events are processed.
//
// # Instrument profiles
//
// Package ipf reads instrument profiles, the descriptions of the instruments (type, exchange, options and so on),
// from IPF files and services; their symbols are the symbols to subscribe to.
//
// # Threads and listeners
//
// The listeners (events, connection state, connectors) are called from the threads of the SDK, or by the executor set
// with [EndpointBuilder.WithExecutor] (an [InPlaceExecutor] runs them in the goroutine that calls
// [InPlaceExecutor.ProcessAllPendingTasks]), not from the goroutine that subscribed. A listener must return quickly:
// it holds back the delivery of the next events, so it must not run long operations or wait for the endpoint, as
// [DXEndpoint.AwaitProcessed] does. The methods documented as callable from a listener, such as adding or removing
// listeners or closing the endpoint, may be called. A panic in a listener is written to the log (package logging) and
// does not end the process.
//
// # Errors
//
// Calls report failures as error values. A Java exception thrown in the SDK is a [common.JavaError]; read it with
// errors.As. The other errors describe the failed native call, for example an isolate that cannot start.
//
// # Related packages
//
//   - logging: the levels, files and a Go listener of the log of the SDK.
//   - Isolate settings: [SetSystemProperty], [GetSystemProperty] and [SetIsolateArguments].
package api
