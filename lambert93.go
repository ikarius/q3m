package q3m

import "math"

// GRS80 ellipsoid constants (identical to WGS84 for practical purposes).
const grs80E = 0.0818191910428158 // first eccentricity

// Lambert93 projection constants.
const lambert93Lambda0 = 3.0 * math.Pi / 180 // central meridian

// Derived constants (precomputed).
const (
	lambert93N  = 0.7256077650532670
	lambert93C  = 11754255.4261
	lambert93Xs = 700000.0
	lambert93Ys = 12655612.0499
)

// isoLat computes the isometric latitude for geodetic latitude phi on
// an ellipsoid with first eccentricity e.
func isoLat(phi, e float64) float64 {
	sinPhi := math.Sin(phi)
	return math.Log(math.Tan(math.Pi/4+phi/2)) - e/2*math.Log((1+e*sinPhi)/(1-e*sinPhi))
}

// ToLambert93 converts WGS84 (lat, lon in degrees) to Lambert93 (E, N in metres).
func ToLambert93(lat, lon float64) (E, N float64) {
	phi := lat * math.Pi / 180
	lambda := lon * math.Pi / 180

	L := isoLat(phi, grs80E)

	r := lambert93C * math.Exp(-lambert93N*L)
	gamma := lambert93N * (lambda - lambert93Lambda0)

	E = lambert93Xs + r*math.Sin(gamma)
	N = lambert93Ys - r*math.Cos(gamma)
	return
}

// FromLambert93 converts Lambert93 (E, N in metres) to WGS84 (lat, lon in degrees).
func FromLambert93(E, N float64) (lat, lon float64) {
	dX := E - lambert93Xs
	dY := lambert93Ys - N

	r := math.Sqrt(dX*dX + dY*dY)
	gamma := math.Atan2(dX, dY)

	lambda := lambert93Lambda0 + gamma/lambert93N
	L := -math.Log(r/lambert93C) / lambert93N

	// Iterative inversion of the isometric latitude.
	phi := 2*math.Atan(math.Exp(L)) - math.Pi/2
	for i := 0; i < 10; i++ {
		eSinPhi := grs80E * math.Sin(phi)
		phiNext := 2*math.Atan(math.Pow((1+eSinPhi)/(1-eSinPhi), grs80E/2)*math.Exp(L)) - math.Pi/2
		if math.Abs(phiNext-phi) < 1e-12 {
			break
		}
		phi = phiNext
	}

	lat = phi * 180 / math.Pi
	lon = lambda * 180 / math.Pi
	return
}
