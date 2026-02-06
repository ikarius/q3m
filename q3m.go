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

// w is the dictionary size, used for base conversion.
const w = uint64(DictSize)

// Encode converts WGS84 coordinates to a q3m three-word address.
func Encode(lat, lon float64) (Address, error) {
	e, n := ToLambert93(lat, lon)

	idx, ok := CellIndex(e, n)
	if !ok {
		return Address{}, fmt.Errorf("q3m: coordinates (%f, %f) are outside the Lambert93 grid", lat, lon)
	}

	shuffled := Shuffle(idx)

	w1 := int(shuffled / (w * w))
	w2 := int((shuffled / w) % w)
	w3 := int(shuffled % w)

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

	var indices [3]uint64
	for i, p := range parts {
		idx, ok := IndexOf(p)
		if !ok {
			return Coordinate{}, fmt.Errorf("q3m: unknown word %q", p)
		}
		indices[i] = uint64(idx)
	}

	shuffled := indices[0]*w*w + indices[1]*w + indices[2]
	idx := Unshuffle(shuffled)

	if idx >= TotalCells {
		return Coordinate{}, fmt.Errorf("q3m: address %q maps to invalid cell index", address)
	}

	e, n := CellCenter(idx)
	lat, lon := FromLambert93(e, n)

	return Coordinate{Lat: lat, Lon: lon}, nil
}
