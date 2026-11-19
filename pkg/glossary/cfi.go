package glossary

import "github.com/dxfeed/dxfeed-graal-go-api/internal/native"

// emptyCFICode is the code of the empty CFI, which is the zero value.
const emptyCFICode = "XXXXXX"

// CFI is a Classification of Financial Instruments code (ISO 10962): six letters, the category, the group and four
// attributes; X stands for an unknown value. The zero value is the empty CFI XXXXXX.
type CFI struct {
	code string // "" for XXXXXX, so that the zero value equals the empty CFI
}

func newCFI(code string) CFI {
	if code == emptyCFICode {
		code = ""
	}
	return CFI{code}
}

// ParseCFI parses a CFI code. A code shorter than six letters is padded with X; an empty code gives the empty CFI.
func ParseCFI(code string) (CFI, error) {
	normalized, err := native.NormalizeCFI(code)
	if err != nil {
		return CFI{}, err
	}
	return newCFI(normalized), nil
}

// CFIOfIntCode returns the CFI of its integer code.
func CFIOfIntCode(intCode int32) (CFI, error) {
	code, err := native.CFIFromIntCode(intCode)
	if err != nil {
		return CFI{}, err
	}
	return newCFI(code), nil
}

// Code returns the six letters of the CFI.
func (c CFI) Code() string {
	if c.code == "" {
		return emptyCFICode
	}
	return c.code
}

// String returns the code of the CFI.
func (c CFI) String() string {
	return c.Code()
}

// IntCode returns the CFI as an integer.
func (c CFI) IntCode() (int32, error) {
	return native.CFIIntCode(c.Code())
}

// Category returns the first letter of the CFI, such as E for equities.
func (c CFI) Category() byte {
	return c.Code()[0]
}

// Group returns the second letter of the CFI, the group in the category.
func (c CFI) Group() byte {
	return c.Code()[1]
}

// IsEquity reports whether the category is E, equities.
func (c CFI) IsEquity() bool { return c.Category() == 'E' }

// IsDebtInstrument reports whether the category is D, debt instruments.
func (c CFI) IsDebtInstrument() bool { return c.Category() == 'D' }

// IsEntitlement reports whether the category is R, entitlements (rights).
func (c CFI) IsEntitlement() bool { return c.Category() == 'R' }

// IsOption reports whether the category is O, options.
func (c CFI) IsOption() bool { return c.Category() == 'O' }

// IsFuture reports whether the category is F, futures.
func (c CFI) IsFuture() bool { return c.Category() == 'F' }

// IsOther reports whether the category is M, others (miscellaneous).
func (c CFI) IsOther() bool { return c.Category() == 'M' }

// CFIAttribute is a position of a CFI code, such as the voting right of an equity.
type CFIAttribute struct {
	Name, Description string
}

// CFIValue is the meaning of a letter of a CFI code at its position.
type CFIValue struct {
	Code              rune
	Name, Description string
	Attribute         CFIAttribute
}

// Decipher returns the meanings of the letters of the CFI in their order.
func (c CFI) Decipher() ([]CFIValue, error) {
	fields, err := native.CFIDecipher(c.Code())
	if err != nil {
		return nil, err
	}
	values := make([]CFIValue, 0, len(fields))
	for _, f := range fields {
		values = append(values, CFIValue{Code: f.Code, Name: f.Name, Description: f.Description,
			Attribute: CFIAttribute{Name: f.AttributeName, Description: f.AttributeDescription}})
	}
	return values, nil
}

// Describe returns the meanings of the letters of the CFI as text.
func (c CFI) Describe() (string, error) {
	return native.CFIDescribe(c.Code())
}
