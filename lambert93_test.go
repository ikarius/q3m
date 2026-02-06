package q3m

import (
	"math"
	"testing"
)

// IGN reference points for Lambert93.
var referencePoints = []struct {
	name    string
	lat     float64
	lon     float64
	easting float64
	northing float64
}{
	{"Paris", 48.8566, 2.3522, 652469.02, 6862035.26},
	{"Marseille", 43.2965, 5.3698, 892390.22, 6247035.26},
	{"Brest", 48.3904, -4.4861, 146632.98, 6836262.33},
	{"Strasbourg", 48.5734, 7.7521, 1050362.70, 6840899.65},
}

func TestToLambert93(t *testing.T) {
	for _, ref := range referencePoints {
		e, n := ToLambert93(ref.lat, ref.lon)
		de := math.Abs(e - ref.easting)
		dn := math.Abs(n - ref.northing)
		// Allow up to 1m tolerance (pyproj reference values).
		if de > 1 || dn > 1 {
			t.Errorf("%s: ToLambert93(%f, %f) = (%f, %f), want ~(%f, %f), delta=(%f, %f)",
				ref.name, ref.lat, ref.lon, e, n, ref.easting, ref.northing, de, dn)
		}
	}
}

func TestFromLambert93(t *testing.T) {
	for _, ref := range referencePoints {
		lat, lon := FromLambert93(ref.easting, ref.northing)
		dlat := math.Abs(lat - ref.lat)
		dlon := math.Abs(lon - ref.lon)
		// Allow up to 0.001 degrees (~100m) since reference coords are rounded.
		if dlat > 0.001 || dlon > 0.001 {
			t.Errorf("%s: FromLambert93(%f, %f) = (%f, %f), want ~(%f, %f)",
				ref.name, ref.easting, ref.northing, lat, lon, ref.lat, ref.lon)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	points := []struct {
		lat float64
		lon float64
	}{
		{48.8584, 2.2945},   // Tour Eiffel
		{43.2965, 5.3698},   // Marseille
		{48.3904, -4.4861},  // Brest
		{42.6887, 9.4507},   // Corse (Bastia)
		{46.5, 3.0},         // Centre France
	}

	for _, p := range points {
		e, n := ToLambert93(p.lat, p.lon)
		lat, lon := FromLambert93(e, n)
		dlat := math.Abs(lat - p.lat)
		dlon := math.Abs(lon - p.lon)
		// Round trip should be better than 1mm.
		if dlat > 1e-8 || dlon > 1e-8 {
			t.Errorf("RoundTrip(%f, %f): got (%f, %f), delta=(%e, %e)",
				p.lat, p.lon, lat, lon, dlat, dlon)
		}
	}
}

func BenchmarkToLambert93(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToLambert93(48.8584, 2.2945)
	}
}

func BenchmarkFromLambert93(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FromLambert93(652469.0, 6862035.0)
	}
}
