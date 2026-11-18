package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

// cString is C.CString freed by the returned function.
func cString(s string) (*C.char, func()) {
	ptr := C.CString(s)
	return ptr, func() { C.free(unsafe.Pointer(ptr)) }
}

// takeString copies and releases a string that the SDK returned; NULL is "" and false.
func takeString(thread *isolateThread, s *C.char) (string, bool) {
	if s == nil {
		return "", false
	}
	defer C.dxfg_v4_String_release(thread.ptr, s)
	return C.GoString(s), true
}

func releaseJavaObject(thread *isolateThread, ptr unsafe.Pointer) {
	if ptr != nil {
		C.dxfg_v4_JavaObjectHandler_release(thread.ptr, (*C.dxfg_v4_java_object_handler)(ptr))
	}
}

// AuthTokenFields are the values of a Java AuthToken.
type AuthTokenFields struct {
	Scheme, Value, HTTPAuthorization string
	User, Password                   string
	HasUser                          bool
}

// AuthToken kinds of NewAuthToken.
const (
	AuthTokenValueOf = iota
	AuthTokenBasic
	AuthTokenBasicUser
	AuthTokenBasicOrNull
	AuthTokenBearer
	AuthTokenBearerOrNull
	AuthTokenCustom
)

// NewAuthToken creates a Java AuthToken of the kind from the arguments, reads its values and releases it. It returns
// nil without an error when an "OrNull" kind gives no token.
func NewAuthToken(kind int, args ...string) (*AuthTokenFields, error) {
	var result *AuthTokenFields
	err := executeInIsolateThread(func(thread *isolateThread) error {
		ptrs := make([]*C.char, len(args))
		for i, arg := range args {
			ptr, free := cString(arg)
			defer free()
			ptrs[i] = ptr
		}
		arg := func(i int) *C.char {
			if i < len(ptrs) {
				return ptrs[i]
			}
			return nil
		}
		var token *C.dxfg_v4_auth_token_t
		var status C.int32_t
		switch kind {
		case AuthTokenValueOf:
			status = C.dxfg_v4_AuthToken_valueOf(thread.ptr, arg(0), &token)
		case AuthTokenBasic:
			status = C.dxfg_v4_AuthToken_createBasicToken(thread.ptr, arg(0), &token)
		case AuthTokenBasicUser:
			status = C.dxfg_v4_AuthToken_createBasicTokenWithUserAndPassword(thread.ptr, arg(0), arg(1), &token)
		case AuthTokenBasicOrNull:
			status = C.dxfg_v4_AuthToken_createBasicTokenOrNull(thread.ptr, arg(0), arg(1), &token)
		case AuthTokenBearer:
			status = C.dxfg_v4_AuthToken_createBearerToken(thread.ptr, arg(0), &token)
		case AuthTokenBearerOrNull:
			status = C.dxfg_v4_AuthToken_createBearerTokenOrNull(thread.ptr, arg(0), &token)
		case AuthTokenCustom:
			status = C.dxfg_v4_AuthToken_createCustomToken(thread.ptr, arg(0), arg(1), &token)
		default:
			return fmt.Errorf("native: unknown auth token kind %d", kind)
		}
		if err := checkStatus(thread, status); err != nil {
			return err
		}
		if token == nil {
			return nil
		}
		defer releaseJavaObject(thread, unsafe.Pointer(token))
		fields := &AuthTokenFields{}
		var s *C.char
		if err := checkStatus(thread, C.dxfg_v4_AuthToken_getScheme(thread.ptr, token, &s)); err != nil {
			return err
		}
		fields.Scheme, _ = takeString(thread, s)
		if err := checkStatus(thread, C.dxfg_v4_AuthToken_getValue(thread.ptr, token, &s)); err != nil {
			return err
		}
		fields.Value, _ = takeString(thread, s)
		if err := checkStatus(thread, C.dxfg_v4_AuthToken_getHttpAuthorization(thread.ptr, token, &s)); err != nil {
			return err
		}
		fields.HTTPAuthorization, _ = takeString(thread, s)
		s = nil
		if err := checkStatus(thread, C.dxfg_v4_AuthToken_getUser(thread.ptr, token, &s)); err != nil {
			return err
		}
		fields.User, fields.HasUser = takeString(thread, s)
		s = nil
		if err := checkStatus(thread, C.dxfg_v4_AuthToken_getPassword(thread.ptr, token, &s)); err != nil {
			return err
		}
		fields.Password, _ = takeString(thread, s)
		result = fields
		return nil
	})
	return result, err
}

// ParseTimePeriod returns the milliseconds of a Java TimePeriod parsed from the string.
func ParseTimePeriod(value string) (int64, error) {
	var millis int64
	err := executeInIsolateThread(func(thread *isolateThread) error {
		ptr, free := cString(value)
		defer free()
		var period *C.dxfg_v4_time_period_t
		if err := checkStatus(thread, C.dxfg_v4_TimePeriod_valueOfString(thread.ptr, ptr, &period)); err != nil {
			return err
		}
		defer releaseJavaObject(thread, unsafe.Pointer(period))
		var time C.int64_t
		if err := checkStatus(thread, C.dxfg_v4_TimePeriod_getTime(thread.ptr, period, &time)); err != nil {
			return err
		}
		millis = int64(time)
		return nil
	})
	return millis, err
}

// FormatTimePeriod returns the string of a Java TimePeriod of the milliseconds.
func FormatTimePeriod(millis int64) (string, error) {
	var result string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var period *C.dxfg_v4_time_period_t
		if err := checkStatus(thread, C.dxfg_v4_TimePeriod_valueOf(thread.ptr, C.int64_t(millis), &period)); err != nil {
			return err
		}
		defer releaseJavaObject(thread, unsafe.Pointer(period))
		var s *C.char
		status := C.dxfg_v4_Object_toString(thread.ptr, (*C.dxfg_v4_java_object_handler)(unsafe.Pointer(period)), &s)
		if err := checkStatus(thread, status); err != nil {
			return err
		}
		result, _ = takeString(thread, s)
		return nil
	})
	return result, err
}

// TimeFormat is a Java TimeFormat; its handle is released when it becomes unreachable.
type TimeFormat struct {
	handle Handler
}

func newTimeFormat(ptr *C.dxfg_v4_time_format_t) *TimeFormat {
	f := &TimeFormat{handle: NewJavaHandle(unsafe.Pointer(ptr))}
	runtime.SetFinalizer(f, func(f *TimeFormat) { _ = f.handle.Free() })
	return f
}

func (f *TimeFormat) ptr() *C.dxfg_v4_time_format_t {
	return (*C.dxfg_v4_time_format_t)(f.handle.Ptr())
}

// Time format kinds of NewTimeFormat.
const (
	TimeFormatDefault = iota
	TimeFormatGMT
)

// NewTimeFormat returns TimeFormat.DEFAULT or TimeFormat.GMT.
func NewTimeFormat(kind int) (*TimeFormat, error) {
	var ptr *C.dxfg_v4_time_format_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		if kind == TimeFormatGMT {
			return checkStatus(thread, C.dxfg_v4_TimeFormat_GMT(thread.ptr, &ptr))
		}
		return checkStatus(thread, C.dxfg_v4_TimeFormat_DEFAULT(thread.ptr, &ptr))
	})
	if err != nil {
		return nil, err
	}
	return newTimeFormat(ptr), nil
}

// ErrUnknownTimeZone is returned by NewTimeFormatInZone for a time zone that Java does not know.
var ErrUnknownTimeZone = errors.New("unknown time zone")

// NewTimeFormatInZone returns the TimeFormat of the time zone; an empty ID is the default time zone of Java.
func NewTimeFormatInZone(zoneID string) (*TimeFormat, error) {
	var ptr *C.dxfg_v4_time_format_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var zone *C.dxfg_v4_time_zone_t
		if zoneID == "" {
			if err := checkStatus(thread, C.dxfg_v4_TimeZone_getDefault(thread.ptr, &zone)); err != nil {
				return err
			}
		} else {
			id, free := cString(zoneID)
			defer free()
			if err := checkStatus(thread, C.dxfg_v4_TimeZone_getTimeZone(thread.ptr, id, &zone)); err != nil {
				return err
			}
		}
		defer releaseJavaObject(thread, unsafe.Pointer(zone))
		if zoneID != "" {
			// Java gives GMT for an unknown ID instead of an error.
			var s *C.char
			if err := checkStatus(thread, C.dxfg_v4_TimeZone_getID(thread.ptr, zone, &s)); err != nil {
				return err
			}
			if got, _ := takeString(thread, s); got != zoneID {
				return fmt.Errorf("%w %q", ErrUnknownTimeZone, zoneID)
			}
		}
		return checkStatus(thread, C.dxfg_v4_TimeFormat_getInstance(thread.ptr, zone, &ptr))
	})
	if err != nil {
		return nil, err
	}
	return newTimeFormat(ptr), nil
}

// Derived time format kinds of Derive.
const (
	TimeFormatWithTimeZone = iota
	TimeFormatWithMillis
	TimeFormatFullISO
)

// Derive returns TimeFormat.withTimeZone, withMillis or asFullIso of the format.
func (f *TimeFormat) Derive(kind int) (*TimeFormat, error) {
	var ptr *C.dxfg_v4_time_format_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		switch kind {
		case TimeFormatWithTimeZone:
			return checkStatus(thread, C.dxfg_v4_TimeFormat_withTimeZone(thread.ptr, f.ptr(), &ptr))
		case TimeFormatWithMillis:
			return checkStatus(thread, C.dxfg_v4_TimeFormat_withMillis(thread.ptr, f.ptr(), &ptr))
		default:
			return checkStatus(thread, C.dxfg_v4_TimeFormat_asFullIso(thread.ptr, f.ptr(), &ptr))
		}
	})
	runtime.KeepAlive(f)
	if err != nil {
		return nil, err
	}
	return newTimeFormat(ptr), nil
}

// Format formats the milliseconds since the Unix epoch.
func (f *TimeFormat) Format(millis int64) (string, error) {
	var result string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var s *C.char
		if err := checkStatus(thread, C.dxfg_v4_TimeFormat_format(thread.ptr, f.ptr(), C.int64_t(millis), &s)); err != nil {
			return err
		}
		result, _ = takeString(thread, s)
		return nil
	})
	runtime.KeepAlive(f)
	return result, err
}

// Parse returns the milliseconds since the Unix epoch of the string.
func (f *TimeFormat) Parse(value string) (int64, error) {
	var millis C.int64_t
	err := executeInIsolateThread(func(thread *isolateThread) error {
		ptr, free := cString(value)
		defer free()
		return checkStatus(thread, C.dxfg_v4_TimeFormat_parse(thread.ptr, f.ptr(), ptr, &millis))
	})
	runtime.KeepAlive(f)
	return int64(millis), err
}
