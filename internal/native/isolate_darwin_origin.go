//go:build darwin && origin

package native

/*
#cgo LDFLAGS: -Wl,-rpath,@executable_path/
*/
import "C"
