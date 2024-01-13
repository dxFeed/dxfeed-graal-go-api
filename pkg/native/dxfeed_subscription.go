package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"

type DXFeedSubscription struct {
	ptr *C.dxfg_subscription_t
}
