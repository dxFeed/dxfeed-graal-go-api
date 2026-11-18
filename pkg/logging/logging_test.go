package logging_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/logging"
)

func TestLevelString(t *testing.T) {
	for level, want := range map[logging.Level]string{logging.All: "ALL", logging.Trace: "TRACE",
		logging.Info: "INFO", logging.Error: "ERROR", logging.Off: "OFF", logging.Level(42): "Level(42)"} {
		if got := level.String(); got != want {
			t.Errorf("%d.String() = %q, want %q", int32(level), got, want)
		}
	}
}

func TestInvalidLevel(t *testing.T) {
	if err := logging.SetLevel(logging.Level(-1)); err == nil {
		t.Error("SetLevel(-1) returned no error")
	}
	if err := logging.Init(logging.Off + 1); err == nil {
		t.Error("Init(Off+1) returned no error")
	}
}

type records struct {
	mu   sync.Mutex
	list []logging.Record
}

func (r *records) add(record logging.Record) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.list = append(r.list, record)
}

func (r *records) find(match func(logging.Record) bool) (logging.Record, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, record := range r.list {
		if match(record) {
			return record, true
		}
	}
	return logging.Record{}, false
}

// TestListenerReceivesSdkRecords connects an endpoint to a closed port: the SDK logs the attempts, and the records
// reach the Go listener.
func TestListenerReceivesSdkRecords(t *testing.T) {
	var got records
	if err := logging.SetListener(got.add); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = logging.SetListener(nil) })
	if err := logging.SetLevel(logging.All); err != nil {
		t.Fatal(err)
	}
	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		t.Fatal(err)
	}
	defer endpoint.Close()
	if err := endpoint.Connect("127.0.0.1:1"); err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	sdkRecord := func(r logging.Record) bool {
		return r.LoggerName != logging.GoLoggerName && strings.Contains(r.Message, "127.0.0.1")
	}
	for {
		if record, ok := got.find(sdkRecord); ok {
			if record.Time.IsZero() || record.ThreadName == "" || !strings.Contains(record.Formatted, record.Message) {
				t.Errorf("incomplete record: %+v", record)
			}
			return
		}
		if time.Since(start) > 10*time.Second {
			t.Fatal("no record of the SDK about 127.0.0.1 within 10 s")
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// initFileChild is set in the child process of TestInitFile to the path of the log file: the logging reads its
// properties once per process.
const initFileChild = "DXFEED_GO_TEST_LOG_FILE"

func TestInitFileChild(t *testing.T) {
	file := os.Getenv(initFileChild)
	if file == "" {
		t.Skip("runs only as a child process")
	}
	if err := logging.InitFile(file, logging.Info); err != nil {
		t.Fatal(err)
	}
	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		t.Fatal(err)
	}
	if err := endpoint.Connect("127.0.0.1:1"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	_ = endpoint.Close()
}

func TestInitFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "sdk.log")
	cmd := exec.Command(os.Args[0], "-test.run=^TestInitFileChild$", "-test.v")
	cmd.Env = append(os.Environ(), initFileChild+"="+file)
	out, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(out), "--- PASS: TestInitFileChild") {
		t.Fatalf("child failed: %v\n%s", err, out)
	}
	if strings.Contains(string(out), "127.0.0.1:1") {
		t.Errorf("the SDK logged to the console instead of the file:\n%s", out)
	}
	// The SDK appends the number of the file of the rotation: sdk.log.0.
	files, err := filepath.Glob(file + "*")
	if err != nil {
		t.Fatal(err)
	}
	var content []byte
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err == nil {
			content = append(content, data...)
		}
	}
	if !strings.Contains(string(content), "127.0.0.1") {
		t.Errorf("the log files %v have no record about 127.0.0.1:\n%s", files, content)
	}
}
