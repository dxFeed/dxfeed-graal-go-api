// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package dxtime_test

import (
	"fmt"
	"log"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/dxtime"
)

// The time format of the SDK, here in GMT: plain, with the time zone, with milliseconds and in full ISO form.
func ExampleFormat_Format() {
	moment := time.Date(2024, 5, 1, 10, 30, 0, 123_000_000, time.UTC)
	gmt := dxtime.GMTFormat()
	for _, f := range []*dxtime.Format{gmt, gmt.WithTimeZone(), gmt.WithMillis(), gmt.AsFullISO()} {
		s, err := f.Format(moment)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(s)
	}
	// Output:
	// 20240501-103000
	// 20240501-103000+0000
	// 20240501-103000.123
	// 2024-05-01T10:30:00.123Z
}

// Parse takes the time forms of the SDK, such as dates, ISO times and times with an offset.
func ExampleFormat_Parse() {
	gmt := dxtime.GMTFormat()
	for _, s := range []string{"20240501-103000", "2024-05-01T10:30:00Z", "2024-05-01T06:30:00-04:00"} {
		t, err := gmt.Parse(s)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(t.UTC().Format(time.RFC3339))
	}
	// Output:
	// 2024-05-01T10:30:00Z
	// 2024-05-01T10:30:00Z
	// 2024-05-01T10:30:00Z
}

// A time period of the SDK, such as an aggregation period, in its ISO 8601 form.
func ExampleParsePeriod() {
	d, err := dxtime.ParsePeriod("PT1H30M")
	if err != nil {
		log.Fatal(err)
	}
	s, err := dxtime.FormatPeriod(d)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(d, s)
	// Output: 1h30m0s PT1H30M0S
}
