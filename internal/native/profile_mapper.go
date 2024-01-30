package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
	"unsafe"
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
func (m *profileMapper) goProfiles(profileList *C.dxfg_instrument_profile_list) []*events.InstrumentProfile {
	if profileList == nil || profileList.elements == nil || int(profileList.size) == 0 {
		return nil
	}

	size := int(profileList.size)
	list := make([]*events.InstrumentProfile, size)
	elementsSlice := unsafe.Slice(profileList.elements, C.size_t(profileList.size))

	for i, event := range elementsSlice {
		nativeEvent := (*C.dxfg_instrument_profile_t)(unsafe.Pointer(event))
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

		customFieldsSize := int(nativeEvent.custom_fields.size)
		customFieldsList := make([]*string, customFieldsSize)

		customFieldsSlice := unsafe.Slice(nativeEvent.custom_fields.elements, C.size_t(customFieldsSize))
		for j, field := range customFieldsSlice {
			customFieldsList[j] = convertString(field)
		}
		profile.SetCustomFields(customFieldsList)

		list[i] = profile
	}

	return list
}
