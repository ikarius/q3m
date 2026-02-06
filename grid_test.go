package q3m

import (
	"math"
	"testing"
)

func TestGridConstants(t *testing.T) {
	if GridWidth != 1_150_000 {
		t.Errorf("GridWidth = %d, want 1150000", GridWidth)
	}
	if GridHeight != 1_070_000 {
		t.Errorf("GridHeight = %d, want 1070000", GridHeight)
	}
	if TotalCells != 1_230_500_000_000 {
		t.Errorf("TotalCells = %d, want 1230500000000", TotalCells)
	}
}

func TestCellIndexOrigin(t *testing.T) {
	idx, ok := CellIndex(EMin, NMin)
	if !ok || idx != 0 {
		t.Errorf("CellIndex at origin: idx=%d, ok=%v", idx, ok)
	}
}

func TestCellIndexTopRight(t *testing.T) {
	// Just inside the top-right corner.
	idx, ok := CellIndex(EMax-0.5, NMax-0.5)
	if !ok {
		t.Fatal("CellIndex at top-right corner returned false")
	}
	expected := (GridHeight-1)*GridWidth + (GridWidth - 1)
	if idx != expected {
		t.Errorf("CellIndex at top-right: got %d, want %d", idx, expected)
	}
}

func TestCellIndexOutOfBounds(t *testing.T) {
	cases := []struct {
		E, N float64
	}{
		{EMin - 1, NMin},
		{EMax, NMin},
		{EMin, NMin - 1},
		{EMin, NMax},
	}
	for _, c := range cases {
		_, ok := CellIndex(c.E, c.N)
		if ok {
			t.Errorf("CellIndex(%f, %f) should be out of bounds", c.E, c.N)
		}
	}
}

func TestCellCenterRoundTrip(t *testing.T) {
	// Encode a point, get its cell, get the centre, re-encode: same cell.
	E, N := 652469.5, 6862035.5
	idx, ok := CellIndex(E, N)
	if !ok {
		t.Fatal("CellIndex failed")
	}
	cE, cN := CellCenter(idx)
	idx2, ok2 := CellIndex(cE, cN)
	if !ok2 || idx2 != idx {
		t.Errorf("CellCenter round-trip failed: %d != %d", idx, idx2)
	}
	// Centre should be within 0.5m of original.
	if math.Abs(cE-E) > 0.5 || math.Abs(cN-N) > 0.5 {
		t.Errorf("Centre too far: (%f,%f) vs (%f,%f)", cE, cN, E, N)
	}
}

func TestMaxIndex(t *testing.T) {
	// The maximum valid index should be TotalCells - 1.
	maxIdx := TotalCells - 1
	E, N := CellCenter(maxIdx)
	idx, ok := CellIndex(E, N)
	if !ok || idx != maxIdx {
		t.Errorf("Max index round-trip: got %d (ok=%v), want %d", idx, ok, maxIdx)
	}
}
