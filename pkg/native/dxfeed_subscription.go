package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type DXFeedSubscription struct {
	ptr *C.dxfg_subscription_t
}

type dxfg_symbol_t struct {
	t      C.int32_t
	symbol *C.char
}

func (s DXFeedSubscription) AddSymbol(symbol string) {
	thread := attachIsolateThread()
	defer thread.detachIsolateThread()
	ss := &dxfg_symbol_t{}
	ss.t = 0
	ss.symbol = C.CString(symbol)
	C.dxfg_DXFeedSubscription_addSymbol(thread.ptr, s.ptr, (*C.dxfg_symbol_t)(unsafe.Pointer(ss)))
}
