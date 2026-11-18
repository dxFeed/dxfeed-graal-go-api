// Command install-sdk installs the dxFeed Graal Native SDK library that this version of the API needs into a per-user
// directory and prints the settings that let cgo link it and the program find it at run time.
//
// Run it from a module that requires the API:
//
//	go run github.com/dxfeed/dxfeed-graal-go-api/cmd/install-sdk
//
// The archive is downloaded from the SDK releases on GitHub and checked against the SHA-256 checksum pinned in the
// API. Nothing is written into the Go module cache and no administrator rights are needed.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/sdk"
)

func main() {
	dir := flag.String("dir", "", "installation directory (default: a directory in the user cache)")
	printDir := flag.Bool("print-dir", false, "print only the installation directory, for scripts")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: install-sdk [-dir DIR] [-print-dir]\n\n"+
			"Installs dxFeed Graal Native SDK %s for %s/%s and prints the build settings.\n\n",
			sdk.Version, runtime.GOOS, runtime.GOARCH)
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(*dir, *printDir); err != nil {
		fmt.Fprintln(os.Stderr, "install-sdk:", err)
		os.Exit(1)
	}
}

func run(dir string, printDir bool) error {
	platform := runtime.GOOS + "/" + runtime.GOARCH
	archive, ok := sdk.Archives[platform]
	if !ok {
		return fmt.Errorf("dxFeed Graal Native SDK %s has no library for %s", sdk.Version, platform)
	}
	if dir == "" {
		var err error
		if dir, err = sdk.DefaultDir(runtime.GOOS, runtime.GOARCH); err != nil {
			return err
		}
	}
	if !sdk.Installed(dir, archive) {
		url := sdk.DownloadURL(archive)
		fmt.Fprintf(os.Stderr, "Downloading %s\n", url)
		client := &http.Client{Timeout: 10 * time.Minute}
		if err := sdk.Install(client, url, archive, dir); err != nil {
			return err
		}
	}
	if printDir {
		fmt.Println(dir)
		return nil
	}
	fmt.Printf("dxFeed Graal Native SDK %s is installed in %s\n\n", sdk.Version, dir)
	printSettings(dir)
	return nil
}

// printSettings prints the environment that builds and runs a program with the library in dir.
func printSettings(dir string) {
	ldflag := quote("-L" + dir)
	if runtime.GOOS == "windows" {
		fmt.Printf("Build with (PowerShell):\n\n"+
			"  $env:CGO_LDFLAGS = \"%s\"\n\n"+
			"Run with the library on PATH, or copy %s next to the program:\n\n"+
			"  $env:PATH = \"%s;\" + $env:PATH\n\n"+
			"To keep the build setting for all shells: go env -w CGO_LDFLAGS=\"%s\"\n",
			ldflag, sdk.LibraryFile(runtime.GOOS), dir, ldflag)
		return
	}
	flags := ldflag + " " + quote("-Wl,-rpath,"+dir)
	fmt.Printf("Build with:\n\n"+
		"  export CGO_LDFLAGS=\"%s\"\n\n"+
		"The run path makes the program find the library in this directory.\n"+
		"To keep the setting for all shells: go env -w CGO_LDFLAGS=\"%s\"\n",
		flags, flags)
}

// quote puts a flag that contains spaces in single quotes, which cgo accepts in CGO_LDFLAGS.
func quote(s string) string {
	if strings.ContainsAny(s, " \t") {
		return "'" + s + "'"
	}
	return s
}
