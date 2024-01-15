//go:build linux

package native

/*
#cgo LDFLAGS: -Wl,-rpath,$ORIGIN -lDxFeedGraalNativeSdk
*/
import "C"
