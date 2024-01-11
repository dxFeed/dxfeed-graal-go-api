package main

import (
	"dxfeed-graal-go-api/pkg/events/market"
	"dxfeed-graal-go-api/pkg/formatutil"
	"fmt"
)

func main() {
	fmt.Println(formatutil.FormatChar('A')) // Output: A
	fmt.Println(formatutil.FormatChar('\r'))
	fmt.Println(formatutil.FormatChar('\000'))

	q := market.NewQuote("AAPL")
	q.SetAskTime(1704993547000)
	fmt.Println(&q)
	//fmt.Println((utils.GetMillisFromTime(1000)))
	//Run(os.Args)
}
