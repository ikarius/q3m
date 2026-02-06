package q3m

import (
	"math"
	"testing"
)

func TestEncodeDecodeTourEiffel(t *testing.T) {
	lat, lon := 48.8584, 2.2945
	addr, err := Encode(lat, lon)
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}
	t.Logf("Tour Eiffel: %s", addr)

	coord, err := Decode(addr.String())
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	// Should be within 0.71m (diagonal of 1m cell) ~ roughly 0.00001 degrees.
	if math.Abs(coord.Lat-lat) > 0.00002 || math.Abs(coord.Lon-lon) > 0.00002 {
		t.Errorf("Round-trip: (%f, %f) -> %s -> (%f, %f), delta=(%e, %e)",
			lat, lon, addr, coord.Lat, coord.Lon,
			math.Abs(coord.Lat-lat), math.Abs(coord.Lon-lon))
	}
}

func TestEncodeDecodeKnownPlaces(t *testing.T) {
	places := []struct {
		name     string
		lat, lon float64
	}{
		{"Tour Eiffel", 48.8584, 2.2945},
		{"Notre-Dame", 48.8530, 2.3499},
		{"Marseille Vieux-Port", 43.2951, 5.3743},
		{"Mont Saint-Michel", 48.6361, -1.5115},
		{"Bastia (Corse)", 42.6970, 9.4503},
		{"Strasbourg Cathedrale", 48.5819, 7.7510},
	}

	for _, p := range places {
		addr, err := Encode(p.lat, p.lon)
		if err != nil {
			t.Errorf("%s: Encode: %v", p.name, err)
			continue
		}

		coord, err := Decode(addr.String())
		if err != nil {
			t.Errorf("%s: Decode(%s): %v", p.name, addr, err)
			continue
		}

		// Max error: 1m cell -> ~0.00001 deg lat, ~0.000015 deg lon.
		if math.Abs(coord.Lat-p.lat) > 0.00002 || math.Abs(coord.Lon-p.lon) > 0.00002 {
			t.Errorf("%s: (%f, %f) -> %s -> (%f, %f)",
				p.name, p.lat, p.lon, addr, coord.Lat, coord.Lon)
		}
	}
}

func TestDecodeInvalidFormat(t *testing.T) {
	_, err := Decode("one.two")
	if err == nil {
		t.Error("expected error for two-word address")
	}
}

func TestDecodeUnknownWord(t *testing.T) {
	_, err := Decode("xyzzy.hello.world")
	if err == nil {
		t.Error("expected error for unknown word")
	}
}

func TestEncodeOutOfBounds(t *testing.T) {
	_, err := Encode(0, 0) // Africa
	if err == nil {
		t.Error("expected error for out-of-bounds coordinates")
	}
}

func TestAllWordsDifferent(t *testing.T) {
	// Encode a few adjacent points and ensure addresses differ.
	a1, _ := Encode(48.8584, 2.2945)
	a2, _ := Encode(48.8585, 2.2945)
	if a1.String() == a2.String() {
		t.Error("adjacent points should have different addresses")
	}
}

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Encode(48.8584, 2.2945)
	}
}

func BenchmarkDecode(b *testing.B) {
	addr, _ := Encode(48.8584, 2.2945)
	s := addr.String()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Decode(s)
	}
}
