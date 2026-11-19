package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"sort"
	"unsafe"
)

// The glossary values are kept in Go as their text; every call creates the native value from the text and releases
// it afterwards.

func goStringAndRelease(thread *isolateThread, s *C.char) string {
	if s == nil {
		return ""
	}
	defer C.dxfg_v4_String_release(thread.ptr, s)
	return C.GoString(s)
}

func withPriceIncrements(thread *isolateThread, text string, f func(h *C.dxfg_v4_price_increments_t) error) error {
	ptr, free := cString(text)
	defer free()
	var h *C.dxfg_v4_price_increments_t
	if err := checkStatus(thread, C.dxfg_v4_PriceIncrements_valueOf(thread.ptr, ptr, &h)); err != nil {
		return err
	}
	defer releaseJavaObject(thread, unsafe.Pointer(h))
	return f(h)
}

func priceIncrementsText(thread *isolateThread, h *C.dxfg_v4_price_increments_t) (string, error) {
	var text *C.char
	if err := checkStatus(thread, C.dxfg_v4_PriceIncrements_getText(thread.ptr, h, &text)); err != nil {
		return "", err
	}
	return goStringAndRelease(thread, text), nil
}

// PriceIncrementsFromValues returns the text of the price increments in the internal representation: an increment,
// or increments alternating with the price limits.
func PriceIncrementsFromValues(values []float64) (string, error) {
	var result string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var list C.dxfg_v4_double_list
		if len(values) > 0 {
			list.size = C.int32_t(len(values))
			list.elements = (*C.double)(C.malloc(C.size_t(len(values)) * C.size_t(unsafe.Sizeof(C.double(0)))))
			if list.elements == nil {
				return errors.New("native: out of memory")
			}
			defer C.free(unsafe.Pointer(list.elements))
			copy(unsafe.Slice((*float64)(unsafe.Pointer(list.elements)), len(values)), values)
		}
		var h *C.dxfg_v4_price_increments_t
		if err := checkStatus(thread, C.dxfg_v4_PriceIncrements_valueOfIncrements(thread.ptr, &list, &h)); err != nil {
			return err
		}
		defer releaseJavaObject(thread, unsafe.Pointer(h))
		var err error
		result, err = priceIncrementsText(thread, h)
		return err
	})
	return result, err
}

// PriceIncrementsValues returns the price increments in the internal representation.
func PriceIncrementsValues(text string) ([]float64, error) {
	if err := checkNoNul(text); err != nil {
		return nil, err
	}
	var result []float64
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withPriceIncrements(thread, text, func(h *C.dxfg_v4_price_increments_t) error {
			var list *C.dxfg_v4_double_list
			if err := checkStatus(thread, C.dxfg_v4_PriceIncrements_getPriceIncrements(thread.ptr, h, &list)); err != nil {
				return err
			}
			defer C.dxfg_v4_double_list_release(thread.ptr, list)
			if list != nil && list.size > 0 {
				result = append([]float64(nil), unsafe.Slice((*float64)(unsafe.Pointer(list.elements)), int(list.size))...)
			}
			return nil
		})
	})
	return result, err
}

// PriceIncrementsOp is an operation of PriceIncrementsCall.
type PriceIncrementsOp int

// The operations of PriceIncrementsCall.
const (
	PriceIncrement PriceIncrementsOp = iota
	PriceIncrementForPrice
	PricePrecision
	PricePrecisionForPrice
	RoundPrice
	RoundPriceWithDirection
	RoundPriceWithMode
	IncrementPrice
	IncrementPriceWithStep
)

// PriceIncrementsArgs are the arguments of the operations; each operation uses some of them.
type PriceIncrementsArgs struct {
	Price, Step  float64
	Direction    int32
	RoundingMode int32
}

// PriceIncrementsCall runs an operation on the price increments of the text.
func PriceIncrementsCall(text string, op PriceIncrementsOp, args PriceIncrementsArgs) (float64, error) {
	if err := checkNoNul(text); err != nil {
		return 0, err
	}
	var result float64
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withPriceIncrements(thread, text, func(h *C.dxfg_v4_price_increments_t) error {
			var value C.double
			var precision C.int32_t
			price, direction := C.double(args.Price), C.int32_t(args.Direction)
			var status C.int32_t
			switch op {
			case PriceIncrement:
				status = C.dxfg_v4_PriceIncrements_getPriceIncrement(thread.ptr, h, &value)
			case PriceIncrementForPrice:
				status = C.dxfg_v4_PriceIncrements_getPriceIncrementForPriceAndDirection(thread.ptr, h, price, direction,
					&value)
			case PricePrecision:
				status = C.dxfg_v4_PriceIncrements_getPricePrecision(thread.ptr, h, &precision)
			case PricePrecisionForPrice:
				status = C.dxfg_v4_PriceIncrements_getPricePrecisionForPrice(thread.ptr, h, price, &precision)
			case RoundPrice:
				status = C.dxfg_v4_PriceIncrements_roundPrice(thread.ptr, h, price, &value)
			case RoundPriceWithDirection:
				status = C.dxfg_v4_PriceIncrements_roundPriceWithDirection(thread.ptr, h, price, direction, &value)
			case RoundPriceWithMode:
				status = C.dxfg_v4_PriceIncrements_roundPriceWithRoundingMode(thread.ptr, h, price,
					C.dxfg_v4_rounding_mode_t(args.RoundingMode), &value)
			case IncrementPrice:
				status = C.dxfg_v4_PriceIncrements_incrementPrice(thread.ptr, h, price, direction, &value)
			case IncrementPriceWithStep:
				status = C.dxfg_v4_PriceIncrements_incrementPriceWithStep(thread.ptr, h, price, direction,
					C.double(args.Step), &value)
			default:
				return fmt.Errorf("native: unknown price increments operation %d", op)
			}
			if err := checkStatus(thread, status); err != nil {
				return err
			}
			if op == PricePrecision || op == PricePrecisionForPrice {
				result = float64(precision)
			} else {
				result = float64(value)
			}
			return nil
		})
	})
	return result, err
}

func withCFI(thread *isolateThread, code string, f func(h *C.dxfg_v4_cfi_t) error) error {
	ptr, free := cString(code)
	defer free()
	var h *C.dxfg_v4_cfi_t
	if err := checkStatus(thread, C.dxfg_v4_CFI_valueOf(thread.ptr, ptr, &h)); err != nil {
		return err
	}
	defer releaseJavaObject(thread, unsafe.Pointer(h))
	return f(h)
}

func cfiCode(thread *isolateThread, h *C.dxfg_v4_cfi_t) (string, error) {
	var code *C.char
	if err := checkStatus(thread, C.dxfg_v4_CFI_getCode(thread.ptr, h, &code)); err != nil {
		return "", err
	}
	return goStringAndRelease(thread, code), nil
}

// NormalizeCFI parses a CFI code and returns the code of the CFI.
func NormalizeCFI(code string) (string, error) {
	if err := checkNoNul(code); err != nil {
		return "", err
	}
	var result string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withCFI(thread, code, func(h *C.dxfg_v4_cfi_t) (err error) {
			result, err = cfiCode(thread, h)
			return err
		})
	})
	return result, err
}

// CFIFromIntCode returns the code of the CFI of an integer code.
func CFIFromIntCode(intCode int32) (string, error) {
	var result string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var h *C.dxfg_v4_cfi_t
		if err := checkStatus(thread, C.dxfg_v4_CFI_valueOfIntCode(thread.ptr, C.int32_t(intCode), &h)); err != nil {
			return err
		}
		defer releaseJavaObject(thread, unsafe.Pointer(h))
		var err error
		result, err = cfiCode(thread, h)
		return err
	})
	return result, err
}

// CFIIntCode returns the integer code of a CFI code.
func CFIIntCode(code string) (int32, error) {
	if err := checkNoNul(code); err != nil {
		return 0, err
	}
	var result int32
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withCFI(thread, code, func(h *C.dxfg_v4_cfi_t) error {
			var value C.int32_t
			if err := checkStatus(thread, C.dxfg_v4_CFI_getIntCode(thread.ptr, h, &value)); err != nil {
				return err
			}
			result = int32(value)
			return nil
		})
	})
	return result, err
}

// CFIDescribe returns the description of a CFI code.
func CFIDescribe(code string) (string, error) {
	if err := checkNoNul(code); err != nil {
		return "", err
	}
	var result string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withCFI(thread, code, func(h *C.dxfg_v4_cfi_t) error {
			var text *C.char
			if err := checkStatus(thread, C.dxfg_v4_CFI_describe(thread.ptr, h, &text)); err != nil {
				return err
			}
			result = goStringAndRelease(thread, text)
			return nil
		})
	})
	return result, err
}

// CFIValueFields are the values of a CFI value: a character of the code and its attribute.
type CFIValueFields struct {
	Code                                rune
	Name, Description                   string
	AttributeName, AttributeDescription string
}

// CFIDecipher returns the values of the characters of a CFI code.
func CFIDecipher(code string) ([]CFIValueFields, error) {
	if err := checkNoNul(code); err != nil {
		return nil, err
	}
	var result []CFIValueFields
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withCFI(thread, code, func(h *C.dxfg_v4_cfi_t) error {
			var list *C.dxfg_v4_cfi_value_list
			if err := checkStatus(thread, C.dxfg_v4_CFI_decipher(thread.ptr, h, &list)); err != nil {
				return err
			}
			defer C.dxfg_v4_cfi_value_list_release(thread.ptr, list)
			if list == nil || list.size == 0 {
				return nil
			}
			for _, v := range unsafe.Slice(list.elements, int(list.size)) {
				fields, err := cfiValueFields(thread, v)
				if err != nil {
					return err
				}
				result = append(result, fields)
			}
			return nil
		})
	})
	return result, err
}

func cfiValueFields(thread *isolateThread, v *C.dxfg_v4_cfi_value_t) (CFIValueFields, error) {
	var fields CFIValueFields
	var code C.int16_t
	var name, description *C.char
	if err := checkStatus(thread, C.dxfg_v4_CFI_Value_getCode(thread.ptr, v, &code)); err != nil {
		return fields, err
	}
	fields.Code = rune(uint16(code))
	if err := checkStatus(thread, C.dxfg_v4_CFI_Value_getName(thread.ptr, v, &name)); err != nil {
		return fields, err
	}
	fields.Name = goStringAndRelease(thread, name)
	if err := checkStatus(thread, C.dxfg_v4_CFI_Value_getDescription(thread.ptr, v, &description)); err != nil {
		return fields, err
	}
	fields.Description = goStringAndRelease(thread, description)
	var attribute *C.dxfg_v4_cfi_attribute_t
	if err := checkStatus(thread, C.dxfg_v4_CFI_Value_getAttribute(thread.ptr, v, &attribute)); err != nil {
		return fields, err
	}
	defer releaseJavaObject(thread, unsafe.Pointer(attribute))
	name, description = nil, nil
	if err := checkStatus(thread, C.dxfg_v4_CFI_Attribute_getName(thread.ptr, attribute, &name)); err != nil {
		return fields, err
	}
	fields.AttributeName = goStringAndRelease(thread, name)
	if err := checkStatus(thread, C.dxfg_v4_CFI_Attribute_getDescription(thread.ptr, attribute,
		&description)); err != nil {
		return fields, err
	}
	fields.AttributeDescription = goStringAndRelease(thread, description)
	return fields, nil
}

func withAdditionalUnderlyings(thread *isolateThread, text string,
	f func(h *C.dxfg_v4_additional_underlyings_t) error) error {
	ptr, free := cString(text)
	defer free()
	var h *C.dxfg_v4_additional_underlyings_t
	if err := checkStatus(thread, C.dxfg_v4_AdditionalUnderlyings_valueOf(thread.ptr, ptr, &h)); err != nil {
		return err
	}
	defer releaseJavaObject(thread, unsafe.Pointer(h))
	return f(h)
}

func additionalUnderlyingsText(thread *isolateThread, h *C.dxfg_v4_additional_underlyings_t) (string, error) {
	var text *C.char
	if err := checkStatus(thread, C.dxfg_v4_AdditionalUnderlyings_getText(thread.ptr, h, &text)); err != nil {
		return "", err
	}
	return goStringAndRelease(thread, text), nil
}

// AdditionalUnderlyingsFromMap returns the text of the additional underlyings of a map from a symbol to its SPC.
func AdditionalUnderlyingsFromMap(m map[string]float64) (string, error) {
	keys := make([]string, 0, len(m))
	for k := range m {
		if err := checkNoNul(k); err != nil {
			return "", err
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var result string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var list C.dxfg_v4_string_to_double_map_entry_list
		if len(keys) > 0 {
			size := C.size_t(len(keys)) * C.size_t(unsafe.Sizeof(C.dxfg_v4_string_to_double_map_entry_t{}))
			list.elements = (*C.dxfg_v4_string_to_double_map_entry_t)(C.calloc(1, size))
			if list.elements == nil {
				return errors.New("native: out of memory")
			}
			entries := unsafe.Slice(list.elements, len(keys))
			defer func() {
				for _, e := range entries {
					C.free(unsafe.Pointer(e.key))
				}
				C.free(unsafe.Pointer(list.elements))
			}()
			for i, k := range keys {
				entries[i].key = C.CString(k)
				entries[i].value = C.double(m[k])
			}
			list.size = C.int32_t(len(keys))
		}
		var h *C.dxfg_v4_additional_underlyings_t
		if err := checkStatus(thread, C.dxfg_v4_AdditionalUnderlyings_valueOfMap(thread.ptr, &list, &h)); err != nil {
			return err
		}
		defer releaseJavaObject(thread, unsafe.Pointer(h))
		var err error
		result, err = additionalUnderlyingsText(thread, h)
		return err
	})
	return result, err
}

// AdditionalUnderlyingsMap returns the map from a symbol to its SPC of the additional underlyings.
func AdditionalUnderlyingsMap(text string) (map[string]float64, error) {
	if err := checkNoNul(text); err != nil {
		return nil, err
	}
	result := map[string]float64{}
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withAdditionalUnderlyings(thread, text, func(h *C.dxfg_v4_additional_underlyings_t) error {
			var list *C.dxfg_v4_string_to_double_map_entry_list
			if err := checkStatus(thread, C.dxfg_v4_AdditionalUnderlyings_getMap(thread.ptr, h, &list)); err != nil {
				return err
			}
			defer C.dxfg_v4_string_to_double_map_entry_list_release(thread.ptr, list)
			if list == nil || list.size == 0 {
				return nil
			}
			for _, e := range unsafe.Slice(list.elements, int(list.size)) {
				result[C.GoString(e.key)] = float64(e.value)
			}
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AdditionalUnderlyingsSPC returns the SPC of a symbol in the additional underlyings, or 0.
func AdditionalUnderlyingsSPC(text, symbol string) (float64, error) {
	if err := checkNoNul(text, symbol); err != nil {
		return 0, err
	}
	var result float64
	err := executeInIsolateThread(func(thread *isolateThread) error {
		return withAdditionalUnderlyings(thread, text, func(h *C.dxfg_v4_additional_underlyings_t) error {
			ptr, free := cString(symbol)
			defer free()
			var spc C.double
			if err := checkStatus(thread, C.dxfg_v4_AdditionalUnderlyings_getSPC(thread.ptr, h, ptr, &spc)); err != nil {
				return err
			}
			result = float64(spc)
			return nil
		})
	})
	return result, err
}

// AdditionalUnderlyingsSPCFromText returns the SPC of a symbol in the text of additional underlyings without
// parsing all of it, or 0.
func AdditionalUnderlyingsSPCFromText(text, symbol string) (float64, error) {
	if err := checkNoNul(text, symbol); err != nil {
		return 0, err
	}
	var result float64
	err := executeInIsolateThread(func(thread *isolateThread) error {
		textPtr, freeText := cString(text)
		defer freeText()
		symbolPtr, freeSymbol := cString(symbol)
		defer freeSymbol()
		var spc C.double
		if err := checkStatus(thread, C.dxfg_v4_AdditionalUnderlyings_getSPCFromText(thread.ptr, textPtr, symbolPtr,
			&spc)); err != nil {
			return err
		}
		result = float64(spc)
		return nil
	})
	return result, err
}
