//go:build darwin

package native

/*
#cgo LDFLAGS: -L${SRCDIR}/graal
#cgo LDFLAGS: -Wl,-rpath,${SRCDIR}/graal
#cgo LDFLAGS: -lDxFeedGraalNativeSdk
*/
import "C"
