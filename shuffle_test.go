package q3m

import "testing"

func TestShuffleRoundTrip(t *testing.T) {
	testCases := []uint64{
		0, 1, 2, 100, 1000, 999_999,
		GridWidth, GridWidth + 1,
		TotalCells - 1, TotalCells - 2,
		TotalCells / 2,
		123_456_789_012,
	}
	for _, idx := range testCases {
		shuffled := Shuffle(idx)
		if shuffled >= TotalCells {
			t.Errorf("Shuffle(%d) = %d, out of range", idx, shuffled)
		}
		back := Unshuffle(shuffled)
		if back != idx {
			t.Errorf("Unshuffle(Shuffle(%d)) = %d", idx, back)
		}
	}
}

func TestShuffleBijective(t *testing.T) {
	// Test bijectivity on a small subset by checking no collisions.
	const N = 100_000
	seen := make(map[uint64]uint64, N)
	for i := uint64(0); i < N; i++ {
		s := Shuffle(i)
		if s >= TotalCells {
			t.Fatalf("Shuffle(%d) = %d out of range", i, s)
		}
		if prev, ok := seen[s]; ok {
			t.Fatalf("Collision: Shuffle(%d) = Shuffle(%d) = %d", prev, i, s)
		}
		seen[s] = i
	}
}

func TestShuffleDecorrelation(t *testing.T) {
	// Adjacent inputs should produce very different outputs.
	prev := Shuffle(0)
	for i := uint64(1); i <= 100; i++ {
		curr := Shuffle(i)
		if curr == prev+1 || curr == prev-1 {
			t.Errorf("Shuffle(%d)=%d and Shuffle(%d)=%d are adjacent - poor decorrelation",
				i-1, prev, i, curr)
		}
		prev = curr
	}
}

func BenchmarkShuffle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Shuffle(uint64(i) % TotalCells)
	}
}

func BenchmarkUnshuffle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Unshuffle(uint64(i) % TotalCells)
	}
}
