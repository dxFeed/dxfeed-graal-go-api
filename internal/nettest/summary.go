// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package nettest

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// The statuses of a run in the summary file.
const (
	StatusOK              = "ok"
	StatusStoppedBySignal = "stopped_by_signal"
	StatusError           = "error"
	StatusNoData          = "no_data"
	StatusNanoMissing     = "nano_missing"
)

// Summary is the summary file of a run, schema nettest-summary/1; a nil pointer is null.
type Summary struct {
	Schema         string             `json:"schema"`
	SpecVersion    string             `json:"spec_version"`
	Tool           SummaryTool        `json:"tool"`
	Side           string             `json:"side"`
	Name           string             `json:"name"`
	Args           []string           `json:"args"`
	Config         map[string]any     `json:"config"`
	StartEpochNs   int64              `json:"start_epoch_ns"`
	EndEpochNs     int64              `json:"end_epoch_ns"`
	Status         string             `json:"status"`
	ExitCode       int                `json:"exit_code"`
	Error          *string            `json:"error"`
	Records        int64              `json:"records"`
	Rate           int64              `json:"rate"`
	TimeSource     string             `json:"time_source"`
	TimeResolution string             `json:"time_resolution"`
	Intervals      []SummaryInterval  `json:"intervals"`
	E2ELag         *SummaryLag        `json:"e2e_lag"`
	PublishCall    *SummaryPublish    `json:"publish_call"`
	ScheduleSlip   *SummarySlip       `json:"schedule_slip"`
	SymbolsSeen    *int64             `json:"symbols_seen"`
	Runtime        SummaryRuntime     `json:"runtime"`
	Connectors     []SummaryConnector `json:"connectors"`
}

// SummaryTool describes the tool.
type SummaryTool struct {
	Implementation string       `json:"implementation"`
	Name           string       `json:"name"`
	Version        *string      `json:"version"`
	SDKVersion     *string      `json:"sdk_version"`
	QDVersion      *string      `json:"qd_version"`
	Build          SummaryBuild `json:"build"`
}

// SummaryBuild describes the runtime of the tool: for the Go tool, the build of the SDK library; a value the tool
// cannot tell is null.
type SummaryBuild struct {
	Vendor *string `json:"vendor"`
	March  *string `json:"march"`
	GC     *string `json:"gc"`
}

// SummaryInterval is a statistics interval.
type SummaryInterval struct {
	StartEpochMs      int64  `json:"start_epoch_ms"`
	EndEpochMs        int64  `json:"end_epoch_ms"`
	Phase             string `json:"phase"`
	Records           int64  `json:"records"`
	Rate              int64  `json:"rate"`
	LagMeanNs         *int64 `json:"lag_mean_ns"`
	LagP50Ns          *int64 `json:"lag_p50_ns"`
	LagP99Ns          *int64 `json:"lag_p99_ns"`
	LagMaxNs          *int64 `json:"lag_max_ns"`
	PublishNsPerEvent *int64 `json:"publish_ns_per_event"`
	SlipMaxNs         *int64 `json:"slip_max_ns"`
	SlipEndNs         *int64 `json:"slip_end_ns"`
	Connections       int    `json:"connections"`
}

// SummaryLag is the e2e lag of the consumer after the warm-up.
type SummaryLag struct {
	Count     int64            `json:"count"`
	Negative  int64            `json:"negative"`
	P50Ns     *int64           `json:"p50_ns"`
	P90Ns     *int64           `json:"p90_ns"`
	P99Ns     *int64           `json:"p99_ns"`
	P999Ns    *int64           `json:"p999_ns"`
	P9999Ns   *int64           `json:"p9999_ns"`
	MaxNs     *int64           `json:"max_ns"`
	MeanNs    *int64           `json:"mean_ns"`
	Histogram SummaryHistogram `json:"histogram"`
}

// SummaryHistogram gives the parameters of the histogram.
type SummaryHistogram struct {
	Lowest  int64 `json:"lowest"`
	Highest int64 `json:"highest"`
	Digits  int   `json:"digits"`
}

// SummaryPublish is the time of the fill and the publishing calls of the publisher.
type SummaryPublish struct {
	Calls          int64  `json:"calls"`
	Events         int64  `json:"events"`
	NsTotal        int64  `json:"ns_total"`
	NsPerEventMean *int64 `json:"ns_per_event_mean"`
}

// SummarySlip is the schedule slip of the publisher with -r.
type SummarySlip struct {
	MaxNs        int64  `json:"max_ns"`
	EndNs        int64  `json:"end_ns"`
	GrowthNsPerS *int64 `json:"growth_ns_per_s"`
}

// SummaryRuntime is the runtime at the end of the run.
type SummaryRuntime struct {
	HeapUsedMb      *int64 `json:"heap_used_mb"`
	HeapCommittedMb *int64 `json:"heap_committed_mb"`
	GCCount         *int64 `json:"gc_count"`
	GCTimeMs        *int64 `json:"gc_time_ms"`
	CPUPercentMean  *int64 `json:"cpu_percent_mean"`
}

// SummaryConnector holds the statistics of a connector.
type SummaryConnector struct {
	Address string              `json:"address"`
	V1      *SummaryConnectorV1 `json:"v1"`
	V2      *SummaryConnectorV2 `json:"v2"`
}

// SummaryConnectorV1 are the connector statistics v1.
type SummaryConnectorV1 struct {
	ReadBytes          int64 `json:"read_bytes"`
	WrittenBytes       int64 `json:"written_bytes"`
	ReadDataRecords    int64 `json:"read_data_records"`
	WrittenDataRecords int64 `json:"written_data_records"`
	ActiveConnections  int64 `json:"active_connections"`
	ClosedConnections  int64 `json:"closed_connections"`
}

// SummaryConnectorV2 are the connector statistics v2.
type SummaryConnectorV2 struct {
	ReadSubRecords    int64 `json:"read_sub_records"`
	WrittenSubRecords int64 `json:"written_sub_records"`
	DataReadLagUs     int64 `json:"data_read_lag_us"`
	DataWriteLagUs    int64 `json:"data_write_lag_us"`
	ReadRTTUs         int64 `json:"read_rtt_us"`
	WriteRTTUs        int64 `json:"write_rtt_us"`
}

// Int returns a pointer to the value, for a field that may be null.
func Int(value int64) *int64 { return &value }

// IntIf returns a pointer to the value when ok, otherwise nil.
func IntIf(ok bool, value int64) *int64 {
	if !ok {
		return nil
	}
	return &value
}

// String returns a pointer to the value, nil for an empty one.
func String(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// NewLag returns the e2e lag of the histogram, the count, the negative lags, the sum and the maximum of the lags.
func NewLag(h *Histogram, negative, sum, maxLag int64) *SummaryLag {
	count := h.Count()
	percentile := func(p float64) *int64 { return IntIf(count > 0, h.ValueAtPercentile(p)) }
	return &SummaryLag{
		Count: count, Negative: negative,
		P50Ns: percentile(50), P90Ns: percentile(90), P99Ns: percentile(99), P999Ns: percentile(99.9),
		P9999Ns: percentile(99.99), MaxNs: IntIf(count > 0, maxLag), MeanNs: IntIf(count > 0, sum/max(count, 1)),
		Histogram: SummaryHistogram{HistogramLowest, HistogramHighest, HistogramDigits},
	}
}

// CheckSummaryPath checks that the directory of the summary file exists and can be written.
func CheckSummaryPath(path string) error {
	if path == "" {
		return usage("--summary cannot write to the directory of ''")
	}
	if info, err := os.Stat(path); err == nil && info.IsDir() || os.IsPathSeparator(path[len(path)-1]) {
		return usage("--summary cannot write to the directory of '%s'", path)
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".nettest-*.tmp")
	if err != nil {
		return usage("--summary cannot write to the directory of '%s'", path)
	}
	name := file.Name()
	_ = file.Close()
	_ = os.Remove(name)
	return nil
}

// WriteSummary writes the summary to <path>.tmp and renames it to the path.
func WriteSummary(path string, s *Summary) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path+".tmp", append(data, '\n'), 0o644); err != nil {
		return err
	}
	if err := os.Rename(path+".tmp", path); err != nil {
		_ = os.Remove(path + ".tmp")
		return err
	}
	return nil
}
