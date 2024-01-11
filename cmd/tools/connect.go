package main

import (
	"dxfeed-graal-go-api/pkg/api/dxendpoint"
	"fmt"
)

type Connect struct{}

func (c Connect) Run(args []string) {
	address := args[2]
	endpoint := dxendpoint.Create()
	endpoint.Connect(address)
	fmt.Print("run connect tool")
}
