package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"math"
	"sort"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

// checkProfiles returns an error for a nil profile or a string with a NUL character, which C would cut.
func checkProfiles(profiles []*events.InstrumentProfile) error {
	for _, p := range profiles {
		if p == nil {
			return errors.New("native: nil instrument profile")
		}
		for _, v := range []int64{p.Icb(), p.Sic(), p.Expiration(), p.LastTrade()} {
			if v < math.MinInt32 || v > math.MaxInt32 {
				return fmt.Errorf("native: instrument profile %q: ICB, SIC, EXPIRATION or LAST_TRADE %d is out of the "+
					"int32 range", symbolOf(p), v)
			}
		}
		for _, s := range profileStrings(p) {
			if s != nil {
				if err := checkNoNul(*s); err != nil {
					return err
				}
			}
		}
		for name, value := range p.CustomFields() {
			if err := checkNoNul(name, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func profileStrings(p *events.InstrumentProfile) []*string {
	return []*string{p.InstrumentType(), p.Symbol(), p.Description(), p.LocalSymbol(), p.LocalDescription(),
		p.Country(), p.Opol(), p.ExchangeData(), p.Exchanges(), p.Currency(), p.BaseCurrency(), p.Cfi(), p.Isin(),
		p.Sedol(), p.Cusip(), p.Product(), p.Underlying(), p.AdditionalUnderlyings(), p.Mmy(), p.OptionType(),
		p.ExpirationStyle(), p.SettlementStyle(), p.PriceIncrements(), p.TradingHours()}
}

// cOptString copies a string to C memory; nil gives NULL.
func cOptString(s *string) *C.char {
	if s == nil {
		return nil
	}
	return C.CString(*s)
}

// newCProfileList builds a list of profile structures in C memory to pass to the SDK; free it with
// freeCProfileList. Check the profiles with checkProfiles first.
func newCProfileList(thread *isolateThread, profiles []*events.InstrumentProfile) (
	*C.dxfg_v4_instrument_profile_data_list, error,
) {
	list := (*C.dxfg_v4_instrument_profile_data_list)(
		C.calloc(1, C.size_t(unsafe.Sizeof(C.dxfg_v4_instrument_profile_data_list{}))))
	if len(profiles) == 0 {
		return list, nil
	}
	var element *C.dxfg_v4_instrument_profile_data_t
	elements := (**C.dxfg_v4_instrument_profile_data_t)(
		C.calloc(C.size_t(len(profiles)), C.size_t(unsafe.Sizeof(element))))
	list.elements = elements
	slice := unsafe.Slice(elements, len(profiles))
	for i, p := range profiles {
		data, err := newCProfile(thread, p)
		if data != nil {
			slice[i] = data
			list.size = C.int32_t(i + 1)
		}
		if err != nil {
			freeCProfileList(thread, list)
			return nil, err
		}
	}
	return list, nil
}

// newCProfile builds a profile structure in C memory; free it with freeCProfile. On error the returned structure,
// if any, still has to be freed.
func newCProfile(thread *isolateThread, p *events.InstrumentProfile) (*C.dxfg_v4_instrument_profile_data_t, error) {
	d := (*C.dxfg_v4_instrument_profile_data_t)(C.calloc(1, C.size_t(unsafe.Sizeof(C.dxfg_v4_instrument_profile_data_t{}))))
	d._type = cOptString(p.InstrumentType())
	d.symbol = cOptString(p.Symbol())
	d.description = cOptString(p.Description())
	d.local_symbol = cOptString(p.LocalSymbol())
	d.local_description = cOptString(p.LocalDescription())
	d.country = cOptString(p.Country())
	d.opol = cOptString(p.Opol())
	d.exchange_data = cOptString(p.ExchangeData())
	d.exchanges = cOptString(p.Exchanges())
	d.currency = cOptString(p.Currency())
	d.base_currency = cOptString(p.BaseCurrency())
	d.cfi = cOptString(p.Cfi())
	d.isin = cOptString(p.Isin())
	d.sedol = cOptString(p.Sedol())
	d.cusip = cOptString(p.Cusip())
	d.icb = C.int32_t(p.Icb())
	d.sic = C.int32_t(p.Sic())
	d.multiplier = C.double(p.Multiplier())
	d.product = cOptString(p.Product())
	d.underlying = cOptString(p.Underlying())
	d.spc = C.double(p.Spc())
	d.additional_underlyings = cOptString(p.AdditionalUnderlyings())
	d.mmy = cOptString(p.Mmy())
	d.expiration = C.int32_t(p.Expiration())
	d.last_trade = C.int32_t(p.LastTrade())
	d.strike = C.double(p.Strike())
	d.option_type = cOptString(p.OptionType())
	d.expiration_style = cOptString(p.ExpirationStyle())
	d.settlement_style = cOptString(p.SettlementStyle())
	d.price_increments = cOptString(p.PriceIncrements())
	d.trading_hours = cOptString(p.TradingHours())

	custom := p.CustomFields()
	if len(custom) == 0 {
		return d, nil
	}
	names := make([]string, 0, len(custom))
	for name := range custom {
		names = append(names, name)
	}
	sort.Strings(names)
	values := make([]string, 0, 2*len(names))
	for _, name := range names {
		values = append(values, name, custom[name])
	}
	nameValues := newStringList(values)
	defer freeStringList(nameValues)
	var fields *C.dxfg_v4_instrument_profile_custom_fields_t
	err := checkStatus(thread, C.dxfg_v4_InstrumentProfileCustomFields_newWithNamesAndValues(thread.ptr, nameValues,
		&fields))
	d.instrument_profile_custom_fields = fields
	return d, err
}

// freeCProfile frees a structure built by newCProfile and releases its custom fields.
func freeCProfile(thread *isolateThread, d *C.dxfg_v4_instrument_profile_data_t) {
	if d == nil {
		return
	}
	for _, s := range []*C.char{d._type, d.symbol, d.description, d.local_symbol, d.local_description, d.country,
		d.opol, d.exchange_data, d.exchanges, d.currency, d.base_currency, d.cfi, d.isin, d.sedol, d.cusip, d.product,
		d.underlying, d.additional_underlyings, d.mmy, d.option_type, d.expiration_style, d.settlement_style,
		d.price_increments, d.trading_hours} {
		C.free(unsafe.Pointer(s))
	}
	releaseJavaObject(thread, unsafe.Pointer(d.instrument_profile_custom_fields))
	C.free(unsafe.Pointer(d))
}

// freeCProfileList frees a list built by newCProfileList, its structures included.
func freeCProfileList(thread *isolateThread, list *C.dxfg_v4_instrument_profile_data_list) {
	if list == nil {
		return
	}
	if list.elements != nil {
		for _, d := range unsafe.Slice(list.elements, int(list.size)) {
			freeCProfile(thread, d)
		}
		C.free(unsafe.Pointer(list.elements))
	}
	C.free(unsafe.Pointer(list))
}

// WriteProfilesToFile writes the profiles to a file; a name ending with .gz or .zip compresses the file.
func WriteProfilesToFile(path string, profiles []*events.InstrumentProfile) error {
	if err := checkNoNul(path); err != nil {
		return err
	}
	if err := checkProfiles(profiles); err != nil {
		return err
	}
	return executeInIsolateThread(func(thread *isolateThread) error {
		list, err := newCProfileList(thread, profiles)
		if err != nil {
			return err
		}
		defer freeCProfileList(thread, list)
		pathPtr, free := cString(path)
		defer free()
		return checkStatus(thread, C.dxfg_v4_InstrumentProfileWriter_writeToFile(thread.ptr, pathPtr, list))
	})
}

// FormatProfileNumber formats a number like the numeric fields of instrument profiles.
func FormatProfileNumber(value float64) (string, error) {
	var result string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var text *C.char
		if err := checkStatus(thread, C.dxfg_v4_InstrumentProfileField_formatNumber(thread.ptr, C.double(value),
			&text)); err != nil {
			return err
		}
		defer C.dxfg_v4_String_release(thread.ptr, text)
		result = C.GoString(text)
		return nil
	})
	return result, err
}

// ParseProfileNumber parses a number like the numeric fields of instrument profiles; an empty text gives 0.
func ParseProfileNumber(text string) (float64, error) {
	if err := checkNoNul(text); err != nil {
		return 0, err
	}
	var result float64
	err := executeInIsolateThread(func(thread *isolateThread) error {
		ptr, free := cString(text)
		defer free()
		var value C.double
		if err := checkStatus(thread, C.dxfg_v4_InstrumentProfileField_parseNumber(thread.ptr, ptr, &value)); err != nil {
			return err
		}
		result = float64(value)
		return nil
	})
	return result, err
}

// FormatProfileDate formats a day id (days since January 1, 1970) as yyyy-MM-dd; 0 gives an empty string.
func FormatProfileDate(dayID int32) (string, error) {
	var result string
	err := executeInIsolateThread(func(thread *isolateThread) error {
		var text *C.char
		if err := checkStatus(thread, C.dxfg_v4_InstrumentProfileField_formatDate(thread.ptr, C.int32_t(dayID),
			&text)); err != nil {
			return err
		}
		defer C.dxfg_v4_String_release(thread.ptr, text)
		result = C.GoString(text)
		return nil
	})
	return result, err
}

// ParseProfileDate parses a yyyy-MM-dd date to a day id; an empty text gives 0.
func ParseProfileDate(text string) (int32, error) {
	if err := checkNoNul(text); err != nil {
		return 0, err
	}
	var result int32
	err := executeInIsolateThread(func(thread *isolateThread) error {
		ptr, free := cString(text)
		defer free()
		var value C.int32_t
		if err := checkStatus(thread, C.dxfg_v4_InstrumentProfileField_parseDate(thread.ptr, ptr, &value)); err != nil {
			return err
		}
		result = int32(value)
		return nil
	})
	return result, err
}

func symbolOf(p *events.InstrumentProfile) string {
	if s := p.Symbol(); s != nil {
		return *s
	}
	return ""
}
