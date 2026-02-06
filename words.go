package q3m

import (
	_ "embed"
	"fmt"
	"strings"
	"sync"
)

//go:embed words_fr.txt
var wordsRaw string

// DictSize is the number of words in the dictionary.
const DictSize = 10800

var (
	wordsOnce  sync.Once
	wordsList  []string
	wordsIndex map[string]int
)

func loadWords() {
	wordsOnce.Do(func() {
		wordsList = strings.Split(strings.TrimSpace(wordsRaw), "\n")
		if len(wordsList) != DictSize {
			panic(fmt.Sprintf("q3m: dictionary has %d words, expected %d", len(wordsList), DictSize))
		}
		wordsIndex = make(map[string]int, DictSize)
		for i, w := range wordsList {
			wordsIndex[w] = i
		}
	})
}

// WordAt returns the word at position i in the dictionary.
func WordAt(i int) string {
	loadWords()
	return wordsList[i]
}

// IndexOf returns the index of word in the dictionary.
// Returns -1 and false if the word is not found.
func IndexOf(word string) (int, bool) {
	loadWords()
	idx, ok := wordsIndex[strings.ToLower(word)]
	if !ok {
		return -1, false
	}
	return idx, true
}
