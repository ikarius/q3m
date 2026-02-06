package q3m

import (
	"fmt"
	"strings"
)

// Coordinate represents a WGS84 position.
type Coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Address represents a q3m three-word address.
type Address struct {
	W1 string `json:"w1"`
	W2 string `json:"w2"`
	W3 string `json:"w3"`
}

// String returns the dotted representation "w1.w2.w3".
func (a Address) String() string {
	return a.W1 + "." + a.W2 + "." + a.W3
}

// W is the dictionary size, used for base conversion.
const W = uint64(DictSize)

// Encode converts WGS84 coordinates to a q3m three-word address.
func Encode(lat, lon float64) (Address, error) {
	e, n := ToLambert93(lat, lon)

	idx, ok := CellIndex(e, n)
	if !ok {
		return Address{}, fmt.Errorf("q3m: coordinates (%f, %f) are outside the Lambert93 grid", lat, lon)
	}

	shuffled := Shuffle(idx)

	w1 := int(shuffled / (W * W))
	w2 := int((shuffled / W) % W)
	w3 := int(shuffled % W)

	return Address{
		W1: WordAt(w1),
		W2: WordAt(w2),
		W3: WordAt(w3),
	}, nil
}

// Decode converts a q3m three-word address (dot-separated) back to WGS84 coordinates.
// The returned coordinate is the centre of the 1m x 1m cell.
func Decode(address string) (Coordinate, error) {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(address)), ".")
	if len(parts) != 3 {
		return Coordinate{}, fmt.Errorf("q3m: invalid address format %q (expected w1.w2.w3)", address)
	}

	i1, ok1 := IndexOf(parts[0])
	i2, ok2 := IndexOf(parts[1])
	i3, ok3 := IndexOf(parts[2])

	if !ok1 {
		return Coordinate{}, fmt.Errorf("q3m: unknown word %q", parts[0])
	}
	if !ok2 {
		return Coordinate{}, fmt.Errorf("q3m: unknown word %q", parts[1])
	}
	if !ok3 {
		return Coordinate{}, fmt.Errorf("q3m: unknown word %q", parts[2])
	}

	shuffled := uint64(i1)*W*W + uint64(i2)*W + uint64(i3)
	idx := Unshuffle(shuffled)

	if idx >= TotalCells {
		return Coordinate{}, fmt.Errorf("q3m: address %q maps to invalid cell index", address)
	}

	e, n := CellCenter(idx)
	lat, lon := FromLambert93(e, n)

	return Coordinate{Lat: lat, Lon: lon}, nil
}
