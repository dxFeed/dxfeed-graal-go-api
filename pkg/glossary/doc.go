// Package glossary holds the values of the instrument profile fields that have a structure of their own: the price
// increments ([PriceIncrements]), the Classification of Financial Instruments code ([CFI]) and the additional
// underlyings of an option ([AdditionalUnderlyings]).
//
// The values are small immutable values that hold their normalized text, so they can be compared with == and kept in
// maps. The methods that need the rules of the SDK compute in the native library and return an error if it fails.
package glossary
