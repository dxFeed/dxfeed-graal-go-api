package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"

func checkCall[T *P, P any](result T) (T, error) {
	return result, nil
}
