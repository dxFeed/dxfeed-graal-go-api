package sdk

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
)

func makeZip(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for name, content := range files {
		f, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func serve(t *testing.T, body []byte) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var requests atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)
	return srv, &requests
}

func checksum(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func TestInstall(t *testing.T) {
	body := makeZip(t, map[string]string{"lib.so": "library", "docs/a.txt": "doc"})
	srv, requests := serve(t, body)
	archive := Archive{Name: "sdk.zip", SHA256: checksum(body)}
	dir := filepath.Join(t.TempDir(), "sdk")

	if err := Install(srv.Client(), srv.URL, archive, dir); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(filepath.Join(dir, "docs", "a.txt")); err != nil || string(data) != "doc" {
		t.Fatalf("docs/a.txt = %q, %v", data, err)
	}
	if !Installed(dir, archive) {
		t.Fatal("Installed = false after Install")
	}

	if err := Install(srv.Client(), srv.URL, archive, dir); err != nil {
		t.Fatal(err)
	}
	if n := requests.Load(); n != 1 {
		t.Errorf("an existing installation was downloaded again: %d requests", n)
	}
	entries, _ := os.ReadDir(filepath.Dir(dir))
	if len(entries) != 1 {
		t.Errorf("temporary files left next to the installation: %v", entries)
	}
}

func TestInstallRefusesChecksumMismatch(t *testing.T) {
	body := makeZip(t, map[string]string{"lib.so": "library"})
	srv, _ := serve(t, body)
	archive := Archive{Name: "sdk.zip", SHA256: strings.Repeat("0", 64)}
	dir := filepath.Join(t.TempDir(), "sdk")

	err := Install(srv.Client(), srv.URL, archive, dir)
	if err == nil || !strings.Contains(err.Error(), "checksum") {
		t.Fatalf("Install error = %v, want a checksum error", err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("the directory was created despite the checksum mismatch: %v", err)
	}
}

func TestInstallRefusesPathOutsideDirectory(t *testing.T) {
	for _, name := range []string{"../evil.txt", "a/../../evil.txt", "/abs.txt", `..\evil.txt`} {
		body := makeZip(t, map[string]string{name: "x"})
		srv, _ := serve(t, body)
		root := t.TempDir()
		dir := filepath.Join(root, "sdk")
		err := Install(srv.Client(), srv.URL, Archive{Name: "sdk.zip", SHA256: checksum(body)}, dir)
		if err == nil || !strings.Contains(err.Error(), "illegal file path") {
			t.Errorf("%q: Install error = %v, want an illegal path error", name, err)
		}
		if _, err := os.Stat(filepath.Join(root, "evil.txt")); err == nil {
			t.Errorf("%q: a file was written outside the directory", name)
		}
	}
}

func TestInstallReplacesOtherContent(t *testing.T) {
	body := makeZip(t, map[string]string{"lib.so": "new"})
	srv, _ := serve(t, body)
	dir := filepath.Join(t.TempDir(), "sdk")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "stale.so"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Install(srv.Client(), srv.URL, Archive{Name: "sdk.zip", SHA256: checksum(body)}, dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "stale.so")); !os.IsNotExist(err) {
		t.Errorf("stale.so survived the installation: %v", err)
	}
}

func TestDefaultDir(t *testing.T) {
	dir, err := DefaultDir("linux", "arm64")
	if err != nil {
		t.Skip(err)
	}
	if !strings.HasSuffix(dir, filepath.Join("graal-native-sdk", Version, "linux-arm64")) {
		t.Errorf("DefaultDir = %s", dir)
	}
}
