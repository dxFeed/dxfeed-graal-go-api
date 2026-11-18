package native

import (
	"log"
	"runtime/debug"
	"sync"
)

// recoverCallback is deferred by every function the SDK calls back. A panic that unwinds into the SDK ends the
// process, so a panic of a user listener is logged with its stack and the callback returns normally.
func recoverCallback(name string) {
	if r := recover(); r != nil {
		log.Printf("dxfeed: %s panicked: %v\n%s", name, r, debug.Stack())
	}
}

// reportedClazzes holds the event classes without a Go mapper that were already reported.
var reportedClazzes sync.Map

// reportUnmappedClazz logs once per event class that its events have no Go mapper and are skipped.
func reportUnmappedClazz(clazz int32) {
	if _, reported := reportedClazzes.LoadOrStore(clazz, true); reported {
		return
	}
	name := "unknown"
	for n, value := range sdkEventClazzes() {
		if value == clazz {
			name = n
			break
		}
	}
	log.Printf("dxfeed: events of class %s (%d) are not supported by the Go API yet and are skipped", name, clazz)
}
