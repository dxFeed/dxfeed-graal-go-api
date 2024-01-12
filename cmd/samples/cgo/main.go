package main

import (
	"dxfeed-graal-go-api/internal/isolate"
	"dxfeed-graal-go-api/pkg/api"
	"math"
	"time"
)

func main() {
	isolate.NewEndpoint(api.Feed).Connect("demo.dxfeed.com:7300")
	time.Sleep(time.Duration(math.MaxInt64))
}
