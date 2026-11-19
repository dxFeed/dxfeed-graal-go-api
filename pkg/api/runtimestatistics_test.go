// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
)

func TestGetRuntimeStatistics(t *testing.T) {
	first, err := api.GetRuntimeStatistics()
	if err != nil {
		t.Fatal(err)
	}
	// Every value is either "not given" (-1) or 0 or more.
	v := reflect.ValueOf(first)
	for i := 0; i < v.NumField(); i++ {
		if n := v.Field(i).Int(); n < -1 {
			t.Errorf("%s = %d", v.Type().Field(i).Name, n)
		}
	}
	// The values that the isolate always has.
	if first.HeapUsedBytes <= 0 || first.ThreadCount <= 0 || first.UptimeMs < 0 || first.AvailableProcessors <= 0 {
		t.Errorf("heap used, threads, uptime, processors: %v", first)
	}
	// Relations between values, where both are given.
	notLess := func(name string, a, b int64) {
		t.Helper()
		if a != -1 && b != -1 && a < b {
			t.Errorf("%s: %d < %d in %v", name, a, b, first)
		}
	}
	notLess("committed heap < used heap", first.HeapCommittedBytes, first.HeapUsedBytes)
	notLess("peak threads < threads", first.PeakThreadCount, first.ThreadCount)
	notLess("started threads < threads", first.TotalStartedThreadCount, first.ThreadCount)

	second, err := api.GetRuntimeStatistics()
	if err != nil {
		t.Fatal(err)
	}
	notLess("uptime went down", second.UptimeMs, first.UptimeMs)
	notLess("CPU time went down", second.ProcessCPUTimeNs, first.ProcessCPUTimeNs)
	notLess("started threads went down", second.TotalStartedThreadCount, first.TotalStartedThreadCount)
}

func TestRuntimeStatisticsString(t *testing.T) {
	s := api.RuntimeStatistics{HeapUsedBytes: 1, YoungGCCount: -1, AvailableProcessors: 18}
	got := s.String()
	for _, want := range []string{"RuntimeStatistics{heapUsedBytes = 1, ", "youngGcCount = -1, ",
		"availableProcessors = 18}"} {
		if !strings.Contains(got, want) {
			t.Errorf("String() = %q, no %q", got, want)
		}
	}
	if n := strings.Count(got, " = "); n != reflect.TypeOf(s).NumField() {
		t.Errorf("String() has %d values, want one per field: %q", n, got)
	}
}
