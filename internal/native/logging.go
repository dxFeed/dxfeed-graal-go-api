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
	LogLevelAll int32 = iota
	LogLevelTrace
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelOff
)

// GoLoggerName is the logger name of the records that the Go API itself reports.
const GoLoggerName = "dxfeed-graal-go-api"

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
		return logListenerErr
	}
	logSink.Store(&fn)
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
	// recoverCallback would report to the same sink: a panicking listener would be called again.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("dxfeed: log listener panicked: %v", r)
		}
	}()
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
	(*sink)(record)
}

// logGo reports a message of the Go API: to the log listener if one is set, otherwise to the standard log.
func logGo(level int32, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	sink := logSink.Load()
	if sink == nil {
		log.Print("dxfeed: " + message)
		return
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("dxfeed: log listener panicked: %v", r)
				log.Print("dxfeed: " + message)
			}
		}()
		(*sink)(LogRecord{
			Level:      level,
			Time:       time.Now(),
			LoggerName: GoLoggerName,
			Message:    message,
			Formatted:  "dxfeed: " + message,
		})
	}()
}
