package main

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"os"
)

func main() {
	// Enable experimental feature, and set scheme for dxLink.
	for _, property := range [][2]string{
		{"dxfeed.experimental.dxlink.enable", "true"},
		{"scheme", "ext:opt:sysprops,resource:dxlink.xml"},
	} {
		if err := api.SetSystemProperty(property[0], property[1]); err != nil {
			fmt.Printf("Error setting %s: %v\n", property[0], err)
			os.Exit(1)
		}
	}
	Run(os.Args)
}
