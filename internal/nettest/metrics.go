// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package nettest

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// MetricsContentType is the content type of the metrics, the text exposition format of Prometheus.
const MetricsContentType = "text/plain; version=0.0.4; charset=utf-8"

// MetricsInfo are the labels of nettest_info.
type MetricsInfo struct {
	Tool, Side, Name, Version, SDKVersion, Role, Event, Collector string
}

// ConnectorMetrics are the cumulative connector statistics of all instances; V2 is nil without them.
type ConnectorMetrics struct {
	ReadBytes, WrittenBytes, ReadRecords, WrittenRecords int64
	V2                                                   *ConnectorMetricsV2
}

// ConnectorMetricsV2 are the sums of the connector statistics v2 in microseconds and the subscription records.
type ConnectorMetricsV2 struct {
	ReadLagUs, RTTUs, SubRecords int64
}

// MetricsSnapshot is a scrape of the metrics; a negative runtime value is not given.
type MetricsSnapshot struct {
	Info MetricsInfo
	// Records and Connections are per instance.
	Records     []int64
	Connections []int
	// Lag is the e2e lag of the consumer, nil on the publisher.
	Lag *LagMetrics
	// The publisher: the time of the fills and the publishing calls, the calls and the slip of the last batch.
	PublishNs, PublishCalls, SlipNs int64
	Connector                       ConnectorMetrics
	Runtime                         RuntimeSample
}

// metricsWriter writes the families of the exposition format.
type metricsWriter struct {
	b      strings.Builder
	labels string
}

func (w *metricsWriter) family(name, kind, help string) {
	fmt.Fprintf(&w.b, "# HELP %s %s\n# TYPE %s %s\n", name, help, name, kind)
}

func (w *metricsWriter) sample(name, extra, value string) {
	labels := w.labels
	if extra != "" {
		labels += "," + extra
	}
	fmt.Fprintf(&w.b, "%s{%s} %s\n", name, labels, value)
}

func (w *metricsWriter) metric(name, kind, help, value string) {
	w.family(name, kind, help)
	if value != "" {
		w.sample(name, "", value)
	}
}

func integer(value int64) string { return strconv.FormatInt(value, 10) }

// seconds formats nanoseconds or microseconds as seconds.
func seconds(value int64, unit time.Duration) string {
	return strconv.FormatFloat(float64(value)*float64(unit)/1e9, 'g', -1, 64)
}

// given formats a runtime value, empty when it is not given.
func given(value int64, format func(int64) string) string {
	if value < 0 {
		return ""
	}
	return format(value)
}

func label(name, value string) string {
	return name + "=" + strconv.Quote(value)
}

// FormatMetrics formats a scrape as the families of the metrics samples of the kit.
func FormatMetrics(s *MetricsSnapshot) string {
	i := s.Info
	w := &metricsWriter{labels: strings.Join([]string{label("tool", i.Tool), label("side", i.Side),
		label("name", i.Name)}, ",")}
	w.family("nettest_info", "gauge", "The configuration of the NetTest run.")
	w.sample("nettest_info", strings.Join([]string{label("spec_version", SpecVersion), label("version", i.Version),
		label("sdk_version", i.SDKVersion), label("role", i.Role), label("event", i.Event),
		label("collector", i.Collector)}, ","), "1")
	consumer := s.Lag != nil
	if consumer {
		w.family("nettest_records_total", "counter", "The received events.")
	} else {
		w.family("nettest_records_total", "counter", "The published events.")
	}
	for n, records := range s.Records {
		w.sample("nettest_records_total", label("connection", strconv.Itoa(n+1)), integer(records))
	}
	if consumer {
		w.family("nettest_e2e_lag_seconds", "histogram", "The receive time minus the time of the event.")
		for n, count := range s.Lag.Buckets {
			le := "+Inf"
			if n < len(LagBuckets) {
				le = seconds(LagBuckets[n], time.Nanosecond)
			}
			w.sample("nettest_e2e_lag_seconds_bucket", label("le", le), integer(count))
		}
		w.sample("nettest_e2e_lag_seconds_sum", "", seconds(s.Lag.SumNs, time.Nanosecond))
		w.sample("nettest_e2e_lag_seconds_count", "", integer(s.Lag.Count))
		w.metric("nettest_e2e_lag_negative_total", "counter", "The events received with a negative lag.",
			integer(s.Lag.Negative))
	} else {
		w.metric("nettest_publish_seconds_total", "counter", "The time of the batch fills and the publishing calls.",
			seconds(s.PublishNs, time.Nanosecond))
		w.metric("nettest_publish_calls_total", "counter", "The publishing calls.", integer(s.PublishCalls))
		w.metric("nettest_schedule_slip_seconds", "gauge", "The schedule slip of the last batch.",
			seconds(s.SlipNs, time.Nanosecond))
	}
	w.family("nettest_connections", "gauge", "The open connections.")
	for n, connections := range s.Connections {
		w.sample("nettest_connections", label("connection", strconv.Itoa(n+1)), strconv.Itoa(connections))
	}
	c := s.Connector
	w.metric("nettest_connector_read_bytes_total", "counter", "The bytes read by the connectors.",
		integer(c.ReadBytes))
	w.metric("nettest_connector_written_bytes_total", "counter", "The bytes written by the connectors.",
		integer(c.WrittenBytes))
	w.metric("nettest_connector_read_records_total", "counter", "The data records read by the connectors.",
		integer(c.ReadRecords))
	w.metric("nettest_connector_written_records_total", "counter", "The data records written by the connectors.",
		integer(c.WrittenRecords))
	if v2 := c.V2; v2 != nil {
		if consumer {
			w.metric("nettest_connector_lag_seconds_total", "counter",
				"The sum of the connector lags of the data records read.", seconds(v2.ReadLagUs, time.Microsecond))
		}
		w.metric("nettest_connector_rtt_seconds_total", "counter",
			"The sum of the round-trip times of the records read and written by the connectors.",
			seconds(v2.RTTUs, time.Microsecond))
		w.metric("nettest_connector_sub_records_total", "counter",
			"The subscription records read and written by the connectors.", integer(v2.SubRecords))
	}
	r := s.Runtime
	w.metric("nettest_runtime_heap_used_bytes", "gauge", "The used heap.", given(r.HeapUsedBytes, integer))
	w.metric("nettest_runtime_heap_committed_bytes", "gauge", "The committed heap.",
		given(r.HeapCommittedBytes, integer))
	w.metric("nettest_runtime_gc_total", "counter", "The collections.", given(r.GCCount, integer))
	w.metric("nettest_runtime_gc_seconds_total", "counter", "The time of the collections.",
		given(r.GCTimeMs, func(ms int64) string { return seconds(ms, time.Millisecond) }))
	w.metric("nettest_runtime_cpu_seconds_total", "counter", "The CPU time of the process.",
		given(r.CPUTimeNs, func(ns int64) string { return seconds(ns, time.Nanosecond) }))
	return w.b.String()
}

// ParseMetricsAddress parses the value of --metrics-port, [<host>:]<port>; the host is 127.0.0.1 by default.
func ParseMetricsAddress(value string) (string, error) {
	host, port := "127.0.0.1", value
	if i := strings.LastIndex(value, ":"); i >= 0 {
		host, port = strings.TrimSuffix(strings.TrimPrefix(value[:i], "["), "]"), value[i+1:]
		if host == "" {
			host = "127.0.0.1"
		}
	}
	number, err := strconv.Atoi(port)
	if err != nil || number < 0 || number > 65535 || strings.TrimLeft(port, "0123456789") != "" {
		return "", usage("--metrics-port needs [<host>:]<port>, where <port> is a number from 0 to 65535: '%s'", value)
	}
	return net.JoinHostPort(host, strconv.Itoa(number)), nil
}

// MetricsServer serves GET /metrics from a function that must not call the native runtime.
type MetricsServer struct {
	listener net.Listener
	server   *http.Server
}

// ListenMetrics binds the address of --metrics-port.
func ListenMetrics(address string) (*MetricsServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &MetricsServer{listener: listener}, nil
}

// URL returns the URL of the metrics, with the bound port.
func (m *MetricsServer) URL() string {
	return "http://" + m.listener.Addr().String() + "/metrics"
}

// Serve serves the scrapes in the background.
func (m *MetricsServer) Serve(scrape func() *MetricsSnapshot) {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", MetricsContentType)
		_, _ = w.Write([]byte(FormatMetrics(scrape())))
	})
	m.server = &http.Server{Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	go func() { _ = m.server.Serve(m.listener) }()
}

// Close stops the server.
func (m *MetricsServer) Close() {
	if m.server != nil {
		_ = m.server.Close()
		return
	}
	_ = m.listener.Close()
}
