package q3m

// feistelKey is the fixed shuffle key (golden ratio * 2^64). IMMUTABLE.
const feistelKey uint64 = 0x9E3779B97F4A7C15

const feistelRounds = 8

// feistelHalfBits is the number of bits for each half.
// TotalCells ~ 1.23e12 fits in 41 bits. We use 21+21=42 bits, with cycle walking.
const feistelHalfBits = 21

const feistelHalfMask = (uint64(1) << feistelHalfBits) - 1 // 0x1FFFFF

// roundFunc is the Feistel round function. It mixes the half-block with a
// round-dependent key using SplitMix64-style mixing.
func roundFunc(val uint64, round uint64) uint64 {
	x := val + round*feistelKey + feistelKey
	x ^= x >> 30
	x *= 0xBF58476D1CE4E5B9
	x ^= x >> 27
	x *= 0x94D049BB133111EB
	x ^= x >> 31
	return x & feistelHalfMask
}

// feistelEncrypt applies a balanced Feistel network permutation.
func feistelEncrypt(val uint64) uint64 {
	left := (val >> feistelHalfBits) & feistelHalfMask
	right := val & feistelHalfMask

	for round := uint64(0); round < feistelRounds; round++ {
		newRight := left ^ roundFunc(right, round)
		left = right
		right = newRight
	}

	return (left << feistelHalfBits) | right
}

// feistelDecrypt inverts feistelEncrypt.
func feistelDecrypt(val uint64) uint64 {
	left := (val >> feistelHalfBits) & feistelHalfMask
	right := val & feistelHalfMask

	for r := feistelRounds - 1; r >= 0; r-- {
		round := uint64(r)
		newLeft := right ^ roundFunc(left, round)
		right = left
		left = newLeft
	}

	return (left << feistelHalfBits) | right
}

// Shuffle applies a bijective permutation on idx within [0, TotalCells).
// Uses cycle walking to handle the non-power-of-2 domain.
func Shuffle(idx uint64) uint64 {
	result := feistelEncrypt(idx)
	for result >= TotalCells {
		result = feistelEncrypt(result)
	}
	return result
}

// Unshuffle inverts Shuffle: given a shuffled index, returns the original.
func Unshuffle(idx uint64) uint64 {
	result := feistelDecrypt(idx)
	for result >= TotalCells {
		result = feistelDecrypt(result)
	}
	return result
}
