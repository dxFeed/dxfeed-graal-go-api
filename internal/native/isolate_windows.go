//go:build windows

package native

// The directory of the library comes from CGO_LDFLAGS (-L<dir>), see cmd/install-sdk;
// at run time Windows finds DxFeedGraalNativeSdk.dll next to the program or on PATH.

/*
#cgo LDFLAGS: -lDxFeedGraalNativeSdk
*/
import "C"
