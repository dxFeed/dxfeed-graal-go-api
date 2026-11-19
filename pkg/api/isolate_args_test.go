// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package api_test

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
)

// isolateArgsChild names the child test that a test of the isolate arguments runs in a new process: the isolate is
// created once per process, so the arguments can be checked only in a process where no call created it yet.
const isolateArgsChild = "DXFEED_GO_TEST_ISOLATE_ARGS"

// TestIsolateArgumentsChild runs only as the child process of the tests below.
func TestIsolateArgumentsChild(t *testing.T) {
	switch os.Getenv(isolateArgsChild) {
	case "":
		t.Skip("runs only as a child process")
	case "applied":
		if err := api.SetIsolateArguments("-Dgo.api.test.first=no"); err != nil {
			t.Fatal(err)
		}
		// A repeated call replaces the arguments.
		if err := api.SetIsolateArguments("-Dgo.api.test.second=a b", "-Xmx256m"); err != nil {
			t.Fatal(err)
		}
		for key, want := range map[string]string{"go.api.test.first": "", "go.api.test.second": "a b"} {
			got, err := api.GetSystemProperty(key)
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Errorf("%s = %q, want %q", key, got, want)
			}
		}
		if err := api.SetIsolateArguments("-Dgo.api.test.third=x"); !errors.Is(err, api.ErrIsolateCreated) {
			t.Errorf("SetIsolateArguments after the creation = %v, want ErrIsolateCreated", err)
		}
	case "rejected":
		if err := api.SetIsolateArguments("-XX:+NoSuchOptionOfTheGoApiTest"); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 2; i++ {
			_, err := api.GetSystemProperty("user.dir")
			if !errors.Is(err, native.ArgumentParsingFailed) {
				t.Errorf("call %d: error = %v, want ArgumentParsingFailed", i, err)
			}
		}
	default:
		t.Fatalf("unknown mode %q", os.Getenv(isolateArgsChild))
	}
}

func runIsolateArgsChild(t *testing.T, mode string) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestIsolateArgumentsChild$", "-test.v")
	cmd.Env = append(os.Environ(), isolateArgsChild+"="+mode)
	out, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(out), "--- PASS: TestIsolateArgumentsChild") {
		t.Fatalf("child %q failed: %v\n%s", mode, err, out)
	}
}

func TestIsolateArgumentsApplied(t *testing.T) {
	runIsolateArgsChild(t, "applied")
}

func TestIsolateArgumentsRejected(t *testing.T) {
	runIsolateArgsChild(t, "rejected")
}

func TestIsolateArgumentsAfterCreation(t *testing.T) {
	if _, err := api.GetSystemProperty("user.dir"); err != nil {
		t.Fatal(err)
	}
	if err := api.SetIsolateArguments("-Dgo.api.test=late"); !errors.Is(err, api.ErrIsolateCreated) {
		t.Errorf("error = %v, want ErrIsolateCreated", err)
	}
}

func TestIsolateArgumentsWithNul(t *testing.T) {
	err := api.SetIsolateArguments("-Dok=1", "-Dbad=a\x00b")
	if err == nil || errors.Is(err, api.ErrIsolateCreated) || !strings.Contains(err.Error(), "#1") {
		t.Errorf("error = %v, want the NUL error of argument #1", err)
	}
}
