//go:build !unix

package native

// handleWriteSignals does nothing: the system has no SIGPIPE or SIGXFSZ.
func handleWriteSignals() {}
