// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

// Package logging configures the logging of the dxFeed Graal SDK and delivers its log records to Go code.
//
// The SDK logs the connections of the endpoints, the errors of the connectors and other events of its Java code. By
// default it writes them to the standard output. There are three ways to change that:
//
//   - [Init], [InitFile] and [InitFiles] set the level and the files of the logging through the Java system properties
//     that the SDK reads when it starts logging. Call them at the start of the program, before the other calls of the
//     API.
//   - [SetListener] gives every record to a Go function instead of the console and the files, for example to pass it
//     to [log/slog]. The records of the Go API itself (a panic of a listener, events that the Go API skips) go to the
//     same function.
//   - [SetLevel] and [SetErrLevel] change the levels at any time.
//
// A Java system property can also be set with the isolate arguments ("-Dlog.level=OFF", see
// [github.com/dxfeed/dxfeed-graal-go-api/pkg/api.SetIsolateArguments]) or in the dxfeed.system.properties file.
package logging

import (
	"fmt"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
)

// Level is a level of the log records, from the most verbose to the least. The SDK uses the levels of
// java.util.logging: FINEST and FINER are given as [Trace], FINE and CONFIG as [Debug], SEVERE as [Error].
type Level int32

// The logging levels.
const (
	// All enables every record.
	All Level = Level(native.LogLevelAll)
	// Trace is a highly detailed tracing message.
	Trace Level = Level(native.LogLevelTrace)
	// Debug is tracing debug information.
	Debug Level = Level(native.LogLevelDebug)
	// Info is an informational message.
	Info Level = Level(native.LogLevelInfo)
	// Warn is a potential problem.
	Warn Level = Level(native.LogLevelWarn)
	// Error is a serious failure.
	Error Level = Level(native.LogLevelError)
	// Off disables the logging.
	Off Level = Level(native.LogLevelOff)
)

var levelNames = [...]string{"ALL", "TRACE", "DEBUG", "INFO", "WARN", "ERROR", "OFF"}

// String returns the name of the level that the system properties log.level and err.level take, for example "INFO".
func (l Level) String() string {
	if l >= All && l <= Off {
		return levelNames[l]
	}
	return fmt.Sprintf("Level(%d)", int32(l))
}

func (l Level) valid() error {
	if l < All || l > Off {
		return fmt.Errorf("logging: invalid level %d", int32(l))
	}
	return nil
}

// GoLoggerName is the [Record.LoggerName] of the records that the Go API itself reports.
const GoLoggerName = native.GoLoggerName

// Record is a log record of the SDK or of the Go API.
type Record struct {
	Level Level
	// Time is the time of the record, with the precision of milliseconds for the records of the SDK.
	Time time.Time
	// ThreadName and ThreadID identify the Java thread that logged the record; empty and 0 for the records of the Go
	// API.
	ThreadName string
	ThreadID   int64
	// LoggerName is the name of the Java logger, or [GoLoggerName].
	LoggerName string
	// Message is the text of the record without the time, the level and the other parts of the layout.
	Message string
	// Err is the Java exception logged with the record (a [github.com/dxfeed/dxfeed-graal-go-api/pkg/common.JavaError]),
	// nil if there is none.
	Err error
	// Formatted is the whole line that the SDK would write: the message with the time, the level, the thread and the
	// exception.
	Formatted string
}

const interceptableLogging = "com.devexperts.logging.InterceptableLogging"

// initLogging sets log.className to InterceptableLogging, as the C++ API does, and then the given properties.
func initLogging(properties ...[2]string) error {
	properties = append([][2]string{{"log.className", interceptableLogging}}, properties...)
	for _, property := range properties {
		if err := native.SetSystemProperty(property[0], property[1]); err != nil {
			return err
		}
	}
	return nil
}

// Init sets the level of the console logging of the SDK. Init(Off) turns the console logging off. It is equivalent to
// setting the system properties log.className to com.devexperts.logging.InterceptableLogging and log.level and
// err.level to the level.
func Init(level Level) error {
	if err := level.valid(); err != nil {
		return err
	}
	return initLogging([2]string{"log.level", level.String()}, [2]string{"err.level", level.String()})
}

// InitFile makes the SDK write its log to the file instead of the console, the records of the level and above. The SDK
// rotates the file and appends the number of the rotation to the name: app.log becomes app.log.0. It is equivalent
// to setting the system properties log.className to com.devexperts.logging.InterceptableLogging, log.file to the file
// and log.level to the level. The properties take effect only before the SDK logs anything, so call it at the start
// of the program.
func InitFile(logFile string, level Level) error {
	if err := level.valid(); err != nil {
		return err
	}
	return initLogging([2]string{"log.file", logFile}, [2]string{"log.level", level.String()})
}

// InitFiles is [InitFile] that also writes the records of errLevel and above (usually [Warn]) to errFile. It is
// equivalent to setting the system properties of InitFile and err.file and err.level.
func InitFiles(logFile, errFile string, logLevel, errLevel Level) error {
	if err := logLevel.valid(); err != nil {
		return err
	}
	if err := errLevel.valid(); err != nil {
		return err
	}
	return initLogging([2]string{"log.file", logFile}, [2]string{"log.level", logLevel.String()},
		[2]string{"err.file", errFile}, [2]string{"err.level", errLevel.String()})
}

// SetLevel sets the level of the records that reach the listener (see [SetListener]), or else the console or the log
// file. The default is [Info] or the system property log.level.
func SetLevel(level Level) error {
	if err := level.valid(); err != nil {
		return err
	}
	return native.SetLogLevel(int32(level))
}

// SetErrLevel sets the level of the err file (see [InitFiles]); it does nothing without one or with a listener.
func SetErrLevel(level Level) error {
	if err := level.valid(); err != nil {
		return err
	}
	return native.SetErrLevel(int32(level))
}

// SetListener gives the log records of the SDK and of the Go API to listener; nil stops the delivery.
//
// The first call with a listener installs it in the SDK for the rest of the process: the SDK stops writing to the
// console and to the files, also after SetListener(nil). Later calls replace the function. The records reach it from
// the thread or goroutine that logs them, one at a time: while it runs, the others that log wait, so it must return
// quickly and must not call the API. The levels ([SetLevel]) filter the records of the SDK only; the few records of
// the Go API always arrive. A panic of the listener is written to the standard log.
func SetListener(listener func(Record)) error {
	if listener == nil {
		return native.SetLogListener(nil)
	}
	return native.SetLogListener(func(r native.LogRecord) {
		listener(Record{
			Level:      Level(r.Level),
			Time:       r.Time,
			ThreadName: r.ThreadName,
			ThreadID:   r.ThreadID,
			LoggerName: r.LoggerName,
			Message:    r.Message,
			Err:        r.Err,
			Formatted:  r.Formatted,
		})
	})
}
