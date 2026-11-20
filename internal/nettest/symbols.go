// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package nettest

// SymbolSeed is the seed of the generated symbols, that of QD NetTest.
const SymbolSeed = 1416948710541751

// javaRandom is java.util.Random as its Javadoc specifies it.
type javaRandom struct{ seed int64 }

const (
	javaMultiplier = 0x5DEECE66D
	javaMask       = 1<<48 - 1
)

func newJavaRandom(seed int64) *javaRandom {
	return &javaRandom{(seed ^ javaMultiplier) & javaMask}
}

func (r *javaRandom) next(bits uint) int32 {
	r.seed = (r.seed*javaMultiplier + 0xB) & javaMask
	return int32(r.seed >> (48 - bits))
}

func (r *javaRandom) nextInt(bound int32) int32 {
	if bound&-bound == bound {
		return int32(int64(bound) * int64(r.next(31)) >> 31)
	}
	for {
		bits := r.next(31)
		value := bits % bound
		if bits-value+(bound-1) >= 0 { // no overflow of int32
			return value
		}
	}
}

// GenerateSymbols returns the first n symbols of QD NetTest in the order of first appearance, and the number of the
// symbols drawn, duplicates included.
func GenerateSymbols(n int) (symbols []string, draws int) {
	r := newJavaRandom(SymbolSeed)
	seen := make(map[string]struct{}, n)
	symbols = make([]string, 0, n)
	buffer := make([]byte, 0, 8)
	for len(symbols) < n {
		length := r.nextInt(3) + r.nextInt(4) + 1
		buffer = buffer[:0]
		for i := int32(0); i < length; i++ {
			if r.nextInt(50) == 0 {
				buffer = append(buffer, byte('0'+r.nextInt(10)))
			} else {
				buffer = append(buffer, byte('A'+r.nextInt(26)))
			}
		}
		draws++
		symbol := string(buffer)
		if _, ok := seen[symbol]; ok {
			continue
		}
		seen[symbol] = struct{}{}
		symbols = append(symbols, symbol)
	}
	return symbols, draws
}

// Slice returns the symbols of the consumer instance (from 1) of instances: those with the index j where
// j % instances == instance - 1.
func Slice(symbols []string, instance, instances int) []string {
	var slice []string
	for j := instance - 1; j < len(symbols); j += instances {
		slice = append(slice, symbols[j])
	}
	return slice
}
