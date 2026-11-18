package logging_test

import (
	"context"
	"log"
	"log/slog"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/logging"
)

// SetListener passes the log of the SDK to log/slog. The listener stays for the rest of the process, so this example
// is not run by the tests.
func ExampleSetListener() {
	if err := logging.SetLevel(logging.Info); err != nil {
		log.Fatal(err)
	}
	err := logging.SetListener(func(r logging.Record) {
		level := slog.LevelInfo
		switch r.Level {
		case logging.Error:
			level = slog.LevelError
		case logging.Warn:
			level = slog.LevelWarn
		case logging.Debug, logging.Trace:
			level = slog.LevelDebug
		}
		attrs := []any{"logger", r.LoggerName, "thread", r.ThreadName}
		if r.Err != nil {
			attrs = append(attrs, "err", r.Err)
		}
		slog.Log(context.Background(), level, r.Message, attrs...)
	})
	if err != nil {
		log.Fatal(err)
	}
}
