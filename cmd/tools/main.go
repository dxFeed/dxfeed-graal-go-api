package main

import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"os"
)

func main() {
	// Enable experimental feature.
	api.SetSystemProperty("dxfeed.experimental.dxlink.enable", "true")
	// Set scheme for dxLink.
	api.SetSystemProperty("scheme", "ext:opt:sysprops,resource:dxlink.xml")
	Run(os.Args)
}
