// Command schedule prints the trading schedule of an instrument at a time: the day and its sessions, the session at
// the time, and the nearest trading days and sessions. The schedule comes from the trading hours of the instrument
// profile of a symbol (the dxFeed demo service is the default source of the profiles) or from a definition:
//
//	go run ./cmd/samples/schedule -symbol AAPL
//	go run ./cmd/samples/schedule -definition "NewYorkETH()" -time 2024-12-24T15:00:00Z
//	go run ./cmd/samples/schedule -defaults schedule.properties -definition "(tz=GMT;td=12345;0=09001700)"
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/ipf"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/schedule"
)

func main() {
	address := flag.String("ipf", "https://demo:demo@tools.dxfeed.com/ipf", "the instrument profile file or URL")
	symbol := flag.String("symbol", "AAPL", "the symbol whose trading hours to use")
	definition := flag.String("definition", "", "the schedule definition to use instead of the instrument profiles")
	defaults := flag.String("defaults", "", "the file of the default schedule data to set first")
	at := flag.String("time", "", "the time in RFC 3339 format, such as 2024-12-24T15:00:00Z; empty for now")
	flag.Parse()

	if err := run(*address, *symbol, *definition, *defaults, *at); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(address, symbol, definition, defaults, at string) error {
	t := time.Now()
	if at != "" {
		var err error
		if t, err = time.Parse(time.RFC3339, at); err != nil {
			return err
		}
	}
	if defaults != "" {
		data, err := os.ReadFile(defaults)
		if err != nil {
			return err
		}
		if err := schedule.SetDefaults(data); err != nil {
			return err
		}
	}
	s, err := scheduleOf(address, symbol, definition)
	if err != nil {
		return err
	}
	zone, err := time.LoadLocation(s.TimeZoneID())
	if err != nil {
		zone = time.UTC
	}
	fmt.Printf("schedule %s, time zone %s\n", s, s.TimeZoneID())

	day, err := s.DayByTime(t)
	if err != nil {
		return err
	}
	fmt.Printf("\n%v, %v, %s - %s\n", day, day.DayOfWeek(), local(day.StartTime(), zone), local(day.EndTime(), zone))
	sessions, err := day.Sessions()
	if err != nil {
		return err
	}
	for _, session := range sessions {
		fmt.Printf("  %-12v %s - %s\n", session.Type(), local(session.StartTime(), zone), local(session.EndTime(), zone))
	}

	session, err := s.SessionByTime(t)
	if err != nil {
		return err
	}
	fmt.Printf("\nat %s: %v session\n", local(t, zone), session.Type())
	if next, ok, err := session.NextSession(schedule.TradingSessions); err != nil {
		return err
	} else if ok {
		fmt.Printf("next trading session: %v %s\n", next.Type(), local(next.StartTime(), zone))
	}
	if prev, ok, err := day.PrevDay(schedule.TradingDays); err != nil {
		return err
	} else if ok {
		fmt.Printf("previous trading day: %v\n", prev)
	}
	if next, ok, err := day.NextDay(schedule.TradingDays); err != nil {
		return err
	} else if ok {
		fmt.Printf("next trading day: %v\n", next)
	}
	if regular, ok, err := s.NearestSessionByTime(t, schedule.RegularSessions); err != nil {
		return err
	} else if ok {
		fmt.Printf("nearest regular session: %s - %s\n", local(regular.StartTime(), zone),
			local(regular.EndTime(), zone))
	}
	return nil
}

func scheduleOf(address, symbol, definition string) (schedule.Schedule, error) {
	if definition != "" {
		return schedule.FromDefinition(definition)
	}
	reader, err := ipf.NewInstrumentProfileReader()
	if err != nil {
		return schedule.Schedule{}, err
	}
	defer func() { _ = reader.Close() }()
	profiles, err := reader.ReadFromFile(address)
	if err != nil {
		return schedule.Schedule{}, err
	}
	for _, p := range profiles {
		if p.Symbol() != nil && *p.Symbol() == symbol {
			return schedule.FromProfile(p)
		}
	}
	return schedule.Schedule{}, fmt.Errorf("no profile of %s in %d profiles", symbol, len(profiles))
}

func local(t time.Time, zone *time.Location) string {
	return t.In(zone).Format("2006-01-02 15:04:05 MST")
}
