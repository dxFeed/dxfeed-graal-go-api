// Command consumer is built by CI as a separate module that uses the API like an application does:
// it creates an endpoint, which loads the native library and checks its version, and closes it.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
)

func main() {
	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		fmt.Fprintln(os.Stderr, "create endpoint:", err)
		os.Exit(1)
	}
	if err := endpoint.CloseAndAwaitTermination(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "close endpoint:", err)
		os.Exit(1)
	}
	fmt.Println("OK: the endpoint was created and closed")
}
