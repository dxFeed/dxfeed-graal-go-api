package ipf

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

func TestReadFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profiles.ipf")
	content := "#STOCK::=TYPE,SYMBOL,DESCRIPTION,EXTRA\n" +
		"STOCK,AAPL,Apple Inc.,first\n" +
		"STOCK,IBM,IBM Corp.,\n" +
		"##COMPLETE\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	reader, err := NewInstrumentProfileReader()
	if err != nil {
		t.Fatalf("new reader: %v", err)
	}
	defer reader.Close()

	profiles, err := reader.ReadFromFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(profiles) != 2 {
		t.Fatalf("read %d profiles, want 2", len(profiles))
	}
	if got := *profiles[0].Symbol(); got != "AAPL" {
		t.Fatalf("symbol = %q, want AAPL", got)
	}
	if got := *profiles[0].Description(); got != "Apple Inc." {
		t.Fatalf("description = %q, want Apple Inc.", got)
	}
	if got, want := profiles[0].CustomFields(), map[string]string{"EXTRA": "first"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("custom fields = %v, want %v", got, want)
	}
	if got := profiles[1].CustomFields(); len(got) != 0 {
		t.Fatalf("custom fields of a profile with an empty field = %v, want none", got)
	}
	complete, err := reader.WasComplete()
	if err != nil || !complete {
		t.Fatalf("WasComplete = %v, %v; want true, nil", complete, err)
	}
	if _, err := reader.LastModified(); err != nil {
		t.Fatalf("LastModified: %v", err)
	}
	url, err := ResolveSourceURL(path)
	if err != nil || url == nil || *url == "" {
		t.Fatalf("ResolveSourceURL = %v, %v", url, err)
	}

	_, err = reader.ReadFromFile(filepath.Join(t.TempDir(), "missing.ipf"))
	var javaErr common.JavaError
	if !errors.As(err, &javaErr) {
		t.Fatalf("error of reading a missing file %v (%T) is not a common.JavaError", err, err)
	}
	if javaErr.ClassName == "" || javaErr.StackTrace == "" {
		t.Fatalf("JavaError without class name or stack trace: %#v", javaErr)
	}
}
