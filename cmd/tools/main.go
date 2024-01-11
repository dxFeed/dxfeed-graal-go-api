package main

import (
	"dxfeed-graal-go-api/pkg/events/market"
	"dxfeed-graal-go-api/pkg/utils"
	"fmt"
)

func main() {
	fmt.Println(utils.FormatChar('A')) // Output: A
	fmt.Println(utils.FormatChar('\r'))
	fmt.Println(utils.FormatChar('\000'))

	q := market.NewQuote("AAPL")
	q.SetAskTime(1704993547000)
	fmt.Println(&q)
	//fmt.Println((utils.GetMillisFromTime(1000)))
	//Run(os.Args)
}
