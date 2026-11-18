package sdk

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// installedMarker is the file written into an installation directory after the archive was verified and extracted;
// it holds the checksum of the archive.
const installedMarker = ".installed"

// LibraryFile returns the file name of the native library on the given operating system (a GOOS value).
func LibraryFile(goos string) string {
	switch goos {
	case "windows":
		return LibraryName + ".dll"
	case "darwin":
		return "lib" + LibraryName + ".dylib"
	default:
		return "lib" + LibraryName + ".so"
	}
}

// DefaultDir returns the per-user directory for the pinned SDK on the given platform:
// <user cache directory>/dxfeed-graal-go-api/graal-native-sdk/<version>/<goos>-<goarch>.
func DefaultDir(goos, goarch string) (string, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("cannot find the user cache directory (choose a directory with -dir): %w", err)
	}
	return filepath.Join(cache, "dxfeed-graal-go-api", "graal-native-sdk", Version, goos+"-"+goarch), nil
}

// Installed reports whether dir holds an installation of the archive.
func Installed(dir string, archive Archive) bool {
	marker, err := os.ReadFile(filepath.Join(dir, installedMarker))
	return err == nil && strings.TrimSpace(string(marker)) == archive.SHA256
}

// Install downloads the archive from url, checks its SHA-256 checksum and extracts it into dir.
// An existing installation of the same archive is kept; any other content of dir is replaced.
// Nothing is written into dir unless the checksum matches.
func Install(client *http.Client, url string, archive Archive, dir string) error {
	if Installed(dir, archive) {
		return nil
	}
	parent := filepath.Dir(dir)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	file, err := os.CreateTemp(parent, ".download-*.zip")
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
		_ = os.Remove(file.Name())
	}()
	if err := download(client, url, file, archive.SHA256); err != nil {
		return err
	}

	staging, err := os.MkdirTemp(parent, ".extract-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(staging) }()
	if err := unzip(file, staging); err != nil {
		return fmt.Errorf("extract %s: %w", archive.Name, err)
	}
	if err := os.WriteFile(filepath.Join(staging, installedMarker), []byte(archive.SHA256+"\n"), 0o644); err != nil {
		return err
	}
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.Rename(staging, dir)
}

// download writes the body of url into file and checks its SHA-256 checksum.
func download(client *http.Client, url string, file *os.File, wantSHA256 string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: %s", url, resp.Status)
	}
	hash := sha256.New()
	if _, err := io.Copy(io.MultiWriter(file, hash), resp.Body); err != nil {
		return fmt.Errorf("download %s: %w", url, err)
	}
	if got := hex.EncodeToString(hash.Sum(nil)); got != wantSHA256 {
		return fmt.Errorf("download %s: SHA-256 checksum is %s, expected %s; the file is not installed",
			url, got, wantSHA256)
	}
	return nil
}

// unzip extracts regular files and directories of the archive into dest; it refuses entries that would land outside
// dest and skips symbolic links and other special files.
func unzip(file *os.File, dest string) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	r, err := zip.NewReader(file, info.Size())
	if err != nil {
		return err
	}
	for _, f := range r.File {
		name := path.Clean(strings.ReplaceAll(f.Name, `\`, "/"))
		if name == "." || path.IsAbs(name) || name == ".." || strings.HasPrefix(name, "../") ||
			filepath.VolumeName(name) != "" {
			return fmt.Errorf("illegal file path in the archive: %q", f.Name)
		}
		target := filepath.Join(dest, filepath.FromSlash(name))
		mode := f.Mode()
		switch {
		case mode.IsDir():
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case mode.IsRegular():
			if err := extractFile(f, target, mode.Perm()|0o644); err != nil {
				return err
			}
		}
	}
	return nil
}

func extractFile(f *zip.File, target string, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	in, err := f.Open()
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	return errors.Join(err, out.Close())
}
