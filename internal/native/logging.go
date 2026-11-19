// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include "dxfg_api.h"
extern void OnLog(graal_isolatethread_t *thread, dxfg_logging_level_t level, int64_t timestamp,
                  char *thread_name, int64_t thread_id, char *logger_name, char *message,
                  dxfg_exception_t *exception, char *formatted_message, void *user_data);
*/
import "C"
import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

// Log levels: the values of dxfg_logging_level_t.
const (
	LogLevelAll   int32 = C.DXFG_LOGGING_LEVEL_ALL
	LogLevelTrace int32 = C.DXFG_LOGGING_LEVEL_TRACE
	LogLevelDebug int32 = C.DXFG_LOGGING_LEVEL_DEBUG
	LogLevelInfo  int32 = C.DXFG_LOGGING_LEVEL_INFO
	LogLevelWarn  int32 = C.DXFG_LOGGING_LEVEL_WARN
	LogLevelError int32 = C.DXFG_LOGGING_LEVEL_ERROR
	LogLevelOff   int32 = C.DXFG_LOGGING_LEVEL_OFF
)

// GoLoggerName is the logger name of the records that the Go API itself reports.
const GoLoggerName = "dxfeed-graal-go-api"

// goLogPrefix starts the messages of the Go API in the standard log.
const goLogPrefix = "dxfeed: "

// LogRecord is a message of the SDK or of the Go API.
type LogRecord struct {
	Level      int32
	Time       time.Time
	ThreadName string
	ThreadID   int64
	LoggerName string
	Message    string
	// Err is the Java exception logged with the message, nil if there is none.
	Err error
	// Formatted is the whole line with the time, the level and the thread; for a record of the Go API it is the
	// line that the standard log would get.
	Formatted string
}

var (
	logSink atomic.Pointer[func(LogRecord)]
	// logDeliveryMu serializes the calls of the sink: the SDK serializes its own records, the records of the Go API
	// come from any goroutine.
	logDeliveryMu sync.Mutex

	logListenerOnce sync.Once
	logListenerErr  error
)

// SetLogListener delivers the log records to fn; nil stops the delivery. The first call with a non-nil fn installs
// the listener of the SDK, which stays installed: the SDK does not log to the console or to files from then on.
func SetLogListener(fn func(LogRecord)) error {
	if fn == nil {
		logSink.Store(nil)
		return nil
	}
	// Stored first, so that the records logged right after the installation are not lost.
	logSink.Store(&fn)
	logListenerOnce.Do(func() {
		logListenerErr = executeInIsolateThread(func(thread *isolateThread) error {
			var listener *C.dxfg_v4_logging_listener_t
			status := C.dxfg_v4_logging_listener_new(thread.ptr, (*[0]byte)(C.OnLog), nil, &listener)
			if err := checkStatus(thread, status); err != nil {
				return err
			}
			// The handle of the listener is kept for the life of the process: releasing it stops the calls.
			return checkStatus(thread, C.dxfg_v4_logging_set_listener(thread.ptr, listener))
		})
	})
	if logListenerErr != nil {
		logSink.Store(nil)
		return logListenerErr
	}
	return nil
}

// SetLogLevel sets the level of the listener, or of the console or of the log file.
func SetLogLevel(level int32) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_logging_set_log_level(thread.ptr, C.dxfg_v4_logging_level_t(level)))
	})
}

// SetErrLevel sets the level of the err file.
func SetErrLevel(level int32) error {
	return executeInIsolateThread(func(thread *isolateThread) error {
		return checkStatus(thread, C.dxfg_v4_logging_set_err_level(thread.ptr, C.dxfg_v4_logging_level_t(level)))
	})
}

// deliverLog gives the record to the sink one call at a time. It returns false if there is no sink or the sink
// panicked; the panic is written to the standard log, not to the sink, which would panic again.
func deliverLog(sink *func(LogRecord), record LogRecord) (delivered bool) {
	if sink == nil {
		return false
	}
	logDeliveryMu.Lock()
	defer logDeliveryMu.Unlock()
	defer func() {
		if r := recover(); r != nil {
			log.Printf(goLogPrefix+"log listener panicked: %v", r)
			delivered = false
		}
	}()
	(*sink)(record)
	return true
}

func goStringOrEmpty(s *C.char) string {
	if s == nil {
		return ""
	}
	return C.GoString(s)
}

//export OnLog
func OnLog(thread *C.graal_isolatethread_t, level C.dxfg_logging_level_t, timestamp C.int64_t, threadName *C.char,
	threadID C.int64_t, loggerName *C.char, message *C.char, exception *C.dxfg_exception_t, formatted *C.char,
	_ unsafe.Pointer) {
	sink := logSink.Load()
	if sink == nil {
		return
	}
	record := LogRecord{
		Level:      int32(level),
		Time:       time.UnixMilli(int64(timestamp)),
		ThreadName: goStringOrEmpty(threadName),
		ThreadID:   int64(threadID),
		LoggerName: goStringOrEmpty(loggerName),
		Message:    goStringOrEmpty(message),
		Formatted:  goStringOrEmpty(formatted),
	}
	if exception != nil {
		record.Err = common.JavaError{
			ClassName:  goStringOrEmpty(exception.class_name),
			Message:    goStringOrEmpty(exception.message),
			StackTrace: goStringOrEmpty(exception.print_stack_trace),
		}
	}
	deliverLog(sink, record)
}

// logGo reports a message of the Go API: to the log listener if one is set, otherwise (or if the listener panics) to
// the standard log.
func logGo(level int32, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	record := LogRecord{
		Level:      level,
		Time:       time.Now(),
		LoggerName: GoLoggerName,
		Message:    message,
		Formatted:  goLogPrefix + message,
	}
	if !deliverLog(logSink.Load(), record) {
		log.Print(record.Formatted)
	}
}
