package main

import (
	"dxfeed-graal-go-api/pkg/api"
	"dxfeed-graal-go-api/pkg/native"
	"math"
	"time"
)

func main() {
	native.NewEndpoint(api.Feed).Connect("demo.dxfeed.com:7300")
	time.Sleep(time.Duration(math.MaxInt64))
}
