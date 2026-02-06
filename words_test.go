package q3m

import "testing"

func TestWordCount(t *testing.T) {
	loadWords()
	if len(wordsList) != DictSize {
		t.Errorf("wordsList has %d entries, want %d", len(wordsList), DictSize)
	}
}

func TestWordAtAndIndexOf(t *testing.T) {
	loadWords()
	for i := 0; i < DictSize; i++ {
		w := WordAt(i)
		idx, ok := IndexOf(w)
		if !ok || idx != i {
			t.Fatalf("IndexOf(WordAt(%d)) = (%d, %v), word=%q", i, idx, ok, w)
		}
	}
}

func TestIndexOfUnknown(t *testing.T) {
	_, ok := IndexOf("xyzzy12345")
	if ok {
		t.Error("IndexOf should return false for unknown word")
	}
}

func TestWordsAreSorted(t *testing.T) {
	loadWords()
	for i := 1; i < len(wordsList); i++ {
		if wordsList[i] <= wordsList[i-1] {
			t.Fatalf("Words not sorted: [%d]=%q >= [%d]=%q", i-1, wordsList[i-1], i, wordsList[i])
		}
	}
}

func TestDictCapacity(t *testing.T) {
	// DictSize^3 must be >= TotalCells.
	capacity := uint64(DictSize) * uint64(DictSize) * uint64(DictSize)
	if capacity < TotalCells {
		t.Errorf("DictSize^3 = %d < TotalCells = %d", capacity, TotalCells)
	}
}

func TestIndexOfCaseInsensitive(t *testing.T) {
	loadWords()
	w := wordsList[0] // first word, e.g. "abandon"
	upper := ""
	for _, r := range w {
		upper += string(r - 32) // ASCII lowercase -> uppercase
	}
	idx1, ok1 := IndexOf(w)
	idx2, ok2 := IndexOf(upper)
	if !ok1 || !ok2 {
		t.Fatalf("IndexOf(%q)=%v, IndexOf(%q)=%v", w, ok1, upper, ok2)
	}
	if idx1 != idx2 {
		t.Errorf("IndexOf(%q)=%d != IndexOf(%q)=%d", w, idx1, upper, idx2)
	}
}

func TestWordsAreASCII(t *testing.T) {
	loadWords()
	for i, w := range wordsList {
		for _, r := range w {
			if r < 'a' || r > 'z' {
				t.Fatalf("word[%d]=%q contains non-ASCII-lowercase rune %q", i, w, r)
			}
		}
	}
}

func TestWordsLength(t *testing.T) {
	loadWords()
	for i, w := range wordsList {
		if len(w) < 4 || len(w) > 8 {
			t.Fatalf("word[%d]=%q has length %d, want 4-8", i, w, len(w))
		}
	}
}

func BenchmarkIndexOf(b *testing.B) {
	loadWords()
	for i := 0; i < b.N; i++ {
		IndexOf(wordsList[i%DictSize])
	}
}
