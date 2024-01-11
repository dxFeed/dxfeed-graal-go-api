//go:build linux

package native

/*
#cgo LDFLAGS: -L${SRCDIR}/graal
#cgo LDFLAGS: -Wl,-rpath,$ORIGIN
#cgo LDFLAGS: -Wl,-rpath,${SRCDIR}/graal
#cgo LDFLAGS: -lDxFeedGraalNativeSdk
*/
import "C"
