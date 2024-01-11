package mathutil

import (
	"reflect"
	"testing"
)

func TestRoundUp(t *testing.T) {
	first := RoundUp(3, 2)
	if first != 4 {
		t.Fatalf(`Not expected value`)
	}
}

func TestCreateEnumBitMaskArrayByValue(t *testing.T) {
	allCases := []int64{1, 2, 3}
	result := CreateEnumBitMaskArrayByValue(-1, allCases)
	expected := []int64{1, 2, 3, -1}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf(`Not expected value %v`, result)
	}
}
