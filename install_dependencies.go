package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	baseUrl := "https://dxfeed.jfrog.io/artifactory/maven/com/dxfeed/graal-native-sdk/%s/graal-native-sdk-%s-%s-%s.zip"
	version := "1.1.6"
	archStr := ""
	osStr := ""
	switch os := runtime.GOOS; os {
	case "darwin":
		osStr = "osx"
	case "linux":
		osStr = "linux"
	case "windows":
		osStr = "windows"
	default:
		panic("Check OS value")
	}
	switch arch := runtime.GOARCH; arch {
	case "arm64":
		archStr = "aarch64"
	case "amd64":
		if osStr == "osx" {
			archStr = "x86_64"
		} else {
			archStr = "amd64"
		}
	default:
		panic("Check Architecture value")
	}
	fullPath := fmt.Sprintf(baseUrl, version, version, archStr, osStr)
	fmt.Printf("Download from %s\n", fullPath)
	fileName := "graal.zip"
	defer os.Remove(fileName)

	err := downloadFile(fileName, fullPath)
	if err != nil {
		panic(err)
	}
	destPath := "internal/native/graal"
	fmt.Printf("Install graal to  %s\n", destPath)

	err = Unzip(fileName, destPath)
	if err != nil {
		panic(err)
	}
}

func downloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
