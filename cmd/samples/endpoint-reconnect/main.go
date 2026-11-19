// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// The endpoint-reconnect sample builds a named endpoint, watches its message connectors and their connection errors,
// and reconnects it on a timer.
//
//	go run ./cmd/samples/endpoint-reconnect -address demo.dxfeed.com:7300 -reconnect 20s -duration 1m
//
// An address where nothing listens, such as 127.0.0.1:1, shows the connection errors.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

// printConnector prints each change of a connector, with the newest error the first time it is seen.
type printConnector struct {
	lastErrorCount int64
}

func (p *printConnector) ConnectorChanged(c api.MessageConnectorInfo) {
	fmt.Printf("%s %s %s: %v, connected to %v\n", time.Now().Format("15:04:05.000"), c.Name, c.Address, c.State,
		c.CurrentAddresses)
	// The listener is called on one SDK thread, one call at a time: no lock is needed.
	if e := c.LastError; e != nil && e.Count != p.lastErrorCount {
		p.lastErrorCount = e.Count
		fmt.Printf("  error #%d at %s: %v\n", e.Count, e.Time.Format("15:04:05.000"), e)
	}
}

type printState func(old, new common.ConnectionState)

func (f printState) UpdateState(old, new common.ConnectionState) { f(old, new) }

func main() {
	address := flag.String("address", "demo.dxfeed.com:7300", "the address to connect to")
	user := flag.String("user", "", "the user name, for the services that require credentials")
	password := flag.String("password", "", "the password, for the services that require credentials")
	reconnect := flag.Duration("reconnect", 20*time.Second, "how often to reconnect")
	duration := flag.Duration("duration", time.Minute, "how long to run")
	flag.Parse()
	if err := run(*address, *user, *password, *reconnect, *duration); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(address, user, password string, reconnect, duration time.Duration) error {
	if reconnect <= 0 {
		return fmt.Errorf("the reconnect period must be positive: %v", reconnect)
	}
	endpoint, err := api.NewEndpointBuilder().
		WithRole(api.Feed).
		WithName("reconnect-sample").
		Build()
	if err != nil {
		return err
	}
	defer func() { _ = endpoint.Close() }()

	endpoint.AddListener(printState(func(old, new common.ConnectionState) {
		fmt.Printf("%s endpoint: %v -> %v\n", time.Now().Format("15:04:05.000"), old, new)
	}))
	if err := endpoint.AddMessageConnectorListener(&printConnector{}); err != nil {
		return err
	}
	if user != "" {
		if err := endpoint.User(user); err != nil {
			return err
		}
		if err := endpoint.Password(password); err != nil {
			return err
		}
	}
	if err := endpoint.Connect(address); err != nil {
		return err
	}

	ticker := time.NewTicker(reconnect)
	defer ticker.Stop()
	stop := time.After(duration)
	for {
		select {
		case <-ticker.C:
			fmt.Println("reconnecting")
			if err := endpoint.Reconnect(); err != nil {
				return err
			}
		case <-stop:
			return printSummary(endpoint)
		}
	}
}

func printSummary(endpoint *api.DXEndpoint) error {
	connectors, err := endpoint.Connectors()
	if err != nil {
		return err
	}
	for _, c := range connectors {
		s := c.Statistics
		fmt.Printf("%s %s: %v, %d connections closed, %d bytes read, %d bytes written\n", c.Name, c.Address, c.State,
			c.ClosedConnectionCount, s.ReadBytes, s.WrittenBytes)
	}
	lastError, err := endpoint.LastConnectionError()
	if err != nil {
		return err
	}
	if lastError == nil {
		fmt.Println("no connection errors")
	} else {
		fmt.Printf("last error (%d in all): %v\n", lastError.Count, lastError)
	}
	return nil
}
