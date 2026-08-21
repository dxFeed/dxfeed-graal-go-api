package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

type profileMapper struct {
}

func newProfileMapper() *profileMapper {
	return &profileMapper{}
}

func convertString(value *C.char) *string {
	if value == nil {
		return nil
	} else {
		result := C.GoString(value)
		return &result
	}
}

func getCustomField(
	thread *isolateThread,
	customFields *C.dxfg_instrument_profile_custom_fields_t,
	name *C.char,
) (string, bool, error) {
	var value *C.char

	err := checkResultCall(func() C.int32_t {
		return C.dxfg_InstrumentProfileCustomFields_getField(thread.ptr,
			customFields,
			name,
			&value)
	})

	if value != nil {
		defer C.dxfg_String_release(thread.ptr, value)
	}
	if err != nil {
		return "", false, err
	}

	if value == nil {
		return "", false, nil
	}

	return C.GoString(value), true, nil
}

func mapCustomFields(
	thread *isolateThread,
	customFields *C.dxfg_instrument_profile_custom_fields_t,
) (map[string]string, error) {
	resultMap := make(map[string]string)

	if customFields == nil {
		return resultMap, nil
	}

	var fieldNames *C.dxfg_string_list

	err := checkResultCall(func() C.int32_t {
		return C.dxfg_InstrumentProfileCustomFields_getNonEmptyFieldNames(thread.ptr,
			customFields,
			&fieldNames)
	})

	if fieldNames != nil {
		defer C.dxfg_CList_String_release(thread.ptr, fieldNames)
	}

	if err != nil {
		return nil, err
	}

	if fieldNames == nil {
		return resultMap, nil
	}

	names := unsafe.Slice(
		fieldNames.elements,
		int(fieldNames.size),
	)

	for _, nativeName := range names {
		if nativeName == nil {
			continue
		}

		value, exists, err := getCustomField(
			thread,
			customFields,
			nativeName,
		)
		if err != nil {
			return nil, err
		}

		if exists {
			resultMap[C.GoString(nativeName)] = value
		}
	}

	return resultMap, nil
}

func (m *profileMapper) goProfiles2(
	thread *isolateThread,
	profileList *C.dxfg_instrument_profile2_list_t,
) ([]*events.InstrumentProfile, error) {
	if profileList == nil ||
		profileList.elements == nil ||
		profileList.size == 0 {
		return nil, nil
	}

	size := int(profileList.size)
	result := make([]*events.InstrumentProfile, 0, size)

	profiles := unsafe.Slice(
		profileList.elements,
		size,
	)

	for _, nativeProfile := range profiles {
		if nativeProfile == nil {
			continue
		}

		nativeEvent := (*C.dxfg_instrument_profile2_t)(unsafe.Pointer(nativeProfile))
		profile := events.NewInstrumentProfile()

		profile.SetSymbol(convertString(nativeEvent.symbol))
		profile.SetInstrumentType(convertString(nativeEvent._type))
		profile.SetDescription(convertString(nativeEvent.description))
		profile.SetLocalSymbol(convertString(nativeEvent.local_symbol))
		profile.SetLocalDescription(convertString(nativeEvent.local_description))
		profile.SetCountry(convertString(nativeEvent.country))
		profile.SetOpol(convertString(nativeEvent.opol))
		profile.SetExchangeData(convertString(nativeEvent.exchange_data))
		profile.SetExchanges(convertString(nativeEvent.exchanges))
		profile.SetCurrency(convertString(nativeEvent.currency))
		profile.SetBaseCurrency(convertString(nativeEvent.base_currency))
		profile.SetCfi(convertString(nativeEvent.cfi))
		profile.SetIsin(convertString(nativeEvent.isin))
		profile.SetSedol(convertString(nativeEvent.sedol))
		profile.SetCusip(convertString(nativeEvent.cusip))
		profile.SetProduct(convertString(nativeEvent.product))
		profile.SetUnderlying(convertString(nativeEvent.underlying))
		profile.SetAdditionalUnderlyings(convertString(nativeEvent.additional_underlyings))
		profile.SetMmy(convertString(nativeEvent.mmy))
		profile.SetOptionType(convertString(nativeEvent.option_type))
		profile.SetExpirationStyle(convertString(nativeEvent.expiration_style))
		profile.SetSettlementStyle(convertString(nativeEvent.settlement_style))
		profile.SetPriceIncrements(convertString(nativeEvent.price_increments))
		profile.SetTradingHours(convertString(nativeEvent.trading_hours))

		profile.SetIcb(int64(nativeEvent.icb))
		profile.SetSic(int64(nativeEvent.sic))
		profile.SetMultiplier(float64(nativeEvent.multiplier))
		profile.SetSpc(float64(nativeEvent.spc))
		profile.SetExpiration(int64(nativeEvent.expiration))
		profile.SetLastTrade(int64(nativeEvent.last_trade))
		profile.SetStrike(float64(nativeEvent.strike))

		customFields, err := mapCustomFields(
			thread,
			nativeEvent.instrument_profile_custom_fields,
		)

		if err != nil {
			return nil, err
		}

		profile.SetCustomFields(customFields)

		result = append(result, profile)
	}

	return result, nil
}
