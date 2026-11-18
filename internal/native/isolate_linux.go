//go:build linux

package native

// The directory of the library comes from CGO_LDFLAGS (-L<dir> -Wl,-rpath,<dir>), see cmd/install-sdk;
// $ORIGIN also lets a program find the library next to it.

/*
#cgo LDFLAGS: -Wl,-rpath,$ORIGIN
#cgo LDFLAGS: -lDxFeedGraalNativeSdk
*/
import "C"
