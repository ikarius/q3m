package q3m

import "math"

// Grid bounds in Lambert93 metres.
const (
	EMin = 100_000.0
	EMax = 1_250_000.0
	NMin = 6_050_000.0
	NMax = 7_120_000.0
)

// Grid dimensions (1m cells).
const (
	GridWidth  uint64 = 1_150_000 // EMax - EMin
	GridHeight uint64 = 1_070_000 // NMax - NMin
	TotalCells uint64 = GridWidth * GridHeight // 1_230_500_000_000
)

// CellIndex returns the grid cell index for the given Lambert93 coordinates.
// Returns false if the point is outside the grid.
func CellIndex(E, N float64) (uint64, bool) {
	if E < EMin || E >= EMax || N < NMin || N >= NMax {
		return 0, false
	}
	x := uint64(math.Floor(E - EMin))
	y := uint64(math.Floor(N - NMin))
	return y*GridWidth + x, true
}

// CellCenter returns the Lambert93 coordinates of the centre of the cell
// identified by idx (+0.5m offset).
func CellCenter(idx uint64) (E, N float64) {
	x := idx % GridWidth
	y := idx / GridWidth
	E = EMin + float64(x) + 0.5
	N = NMin + float64(y) + 0.5
	return
}
