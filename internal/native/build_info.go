// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package native

/*
#include <stdlib.h>
#include "dxfg_api.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

// ErrBuildInfoUnsupported is returned when the library does not report the feature DXFG_V4_FEATURE_BUILD_INFO.
var ErrBuildInfoUnsupported = errors.New("the dxFeed Graal Native SDK library does not give its build information " +
	"(feature " + C.DXFG_V4_FEATURE_BUILD_INFO + ")")

// BuildInfoData holds the values of dxfg_v4_system_get_build_info; a value the library does not know is empty.
type BuildInfoData struct {
	QDVersion string
	Vendor    string
	March     string
	GC        string
}

// GetBuildInfo reads the build information of the library, or returns ErrBuildInfoUnsupported.
func GetBuildInfo() (BuildInfoData, error) {
	supported, err := IsFeatureSupported(C.DXFG_V4_FEATURE_BUILD_INFO)
	if err != nil {
		return BuildInfoData{}, err
	}
	if !supported {
		return BuildInfoData{}, ErrBuildInfoUnsupported
	}
	var d BuildInfoData
	err = executeInIsolateThread(func(thread *isolateThread) error {
		for _, v := range []struct {
			key string
			dst *string
		}{
			{C.DXFG_V4_BUILD_INFO_QD_VERSION, &d.QDVersion},
			{C.DXFG_V4_BUILD_INFO_VENDOR, &d.Vendor},
			{C.DXFG_V4_BUILD_INFO_MARCH, &d.March},
			{C.DXFG_V4_BUILD_INFO_GC, &d.GC},
		} {
			if err := readBuildInfo(thread, v.key, v.dst); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return BuildInfoData{}, err
	}
	return d, nil
}

// readBuildInfo reads one value of the build information into dst; NULL leaves it empty.
func readBuildInfo(thread *isolateThread, key string, dst *string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	var value *C.char
	if err := checkStatus(thread, C.dxfg_v4_system_get_build_info(thread.ptr, cKey, &value)); err != nil {
		return err
	}
	*dst = takeString(thread, value)
	return nil
}
