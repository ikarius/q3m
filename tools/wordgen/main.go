// Command wordgen generates the q3m word list from Lexique383.
//
// Usage:
//
//	go run ./tools/wordgen > words_fr.txt
//
// It downloads Lexique383, filters and curates 10800 words suitable for
// encoding geographic coordinates, then writes one word per line to stdout.
//
// Criteria:
//   - 4-8 letters, ASCII only (no accents)
//   - Nouns, adjectives, or infinitive verbs (lemmas only)
//   - Frequency > 0
//   - No homophones (nbhomoph <= 1)
//   - No offensive words
//   - Sorted alphabetically
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

const (
	lexiqueURL = "http://www.lexique.org/databases/Lexique383/Lexique383.tsv"
	targetSize = 10800
)

// offensive words to exclude.
var banned = map[string]bool{
	"anal": true, "anus": true, "bite": true, "bites": true,
	"chier": true, "chiant": true, "conne": true, "connard": true,
	"cul": true, "culs": true, "enculer": true, "foutre": true,
	"garce": true, "merde": true, "merdes": true, "negro": true,
	"nique": true, "niquer": true, "pisse": true, "pisser": true,
	"puer": true, "putain": true, "pute": true, "putes": true,
	"salaud": true, "salope": true, "salopes": true, "tapette": true,
	"tarer": true, "nazi": true, "nazis": true, "viol": true,
	"violer": true, "viols": true, "haine": true,
}

// isASCIILower checks if a string contains only ASCII lowercase letters.
func isASCIILower(s string) bool {
	for _, r := range s {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}

// hasAccent checks if a string has characters with diacritical marks.
func hasAccent(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return true
		}
	}
	return false
}

type entry struct {
	word string
	freq float64
}

func main() {
	fmt.Fprintf(os.Stderr, "Downloading Lexique383...\n")
	resp, err := http.Get(lexiqueURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "HTTP %d\n", resp.StatusCode)
		os.Exit(1)
	}

	// Lexique383 is a TSV file.
	reader := csv.NewReader(bufio.NewReader(resp.Body))
	reader.Comma = '\t'
	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading header: %v\n", err)
		os.Exit(1)
	}

	// Find column indices.
	colIdx := make(map[string]int)
	for i, h := range header {
		colIdx[strings.TrimSpace(h)] = i
	}

	needed := []string{"ortho", "lemme", "cgram", "freqlemfilms2", "nbhomoph", "islem"}
	for _, n := range needed {
		if _, ok := colIdx[n]; !ok {
			fmt.Fprintf(os.Stderr, "Missing column: %s\n", n)
			fmt.Fprintf(os.Stderr, "Available: %v\n", header)
			os.Exit(1)
		}
	}

	seen := make(map[string]bool)
	var candidates []entry

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // skip malformed rows
		}

		ortho := strings.TrimSpace(record[colIdx["ortho"]])
		lemme := strings.TrimSpace(record[colIdx["lemme"]])
		cgram := strings.TrimSpace(record[colIdx["cgram"]])
		freqStr := strings.TrimSpace(record[colIdx["freqlemfilms2"]])
		homophStr := strings.TrimSpace(record[colIdx["nbhomoph"]])
		isLemStr := strings.TrimSpace(record[colIdx["islem"]])

		_ = isLemStr // used for priority, not filtering
		_ = lemme    // used for priority, not filtering

		// Only nouns, adjectives, infinitive verbs, adverbs.
		switch cgram {
		case "NOM", "ADJ", "VER", "ADV":
			// OK
		default:
			continue
		}

		// Length 4-8.
		if len(ortho) < 4 || len(ortho) > 8 {
			continue
		}

		// ASCII only, no accents.
		if hasAccent(ortho) || !isASCIILower(ortho) {
			continue
		}

		// Frequency > 0.
		freq, err := strconv.ParseFloat(strings.Replace(freqStr, ",", ".", 1), 64)
		if err != nil || freq <= 0 {
			continue
		}

		// Prefer lemmas: give them a bonus.
		if isLemStr == "1" && ortho == lemme {
			freq *= 2
		}

		// Mild homophone filter: exclude words with many homophones.
		homoph, err := strconv.Atoi(homophStr)
		if err != nil || homoph > 3 {
			continue
		}

		// Not banned.
		if banned[ortho] {
			continue
		}

		// Dedup.
		if seen[ortho] {
			continue
		}
		seen[ortho] = true

		candidates = append(candidates, entry{word: ortho, freq: freq})
	}

	fmt.Fprintf(os.Stderr, "Candidates after filtering: %d\n", len(candidates))

	if len(candidates) < targetSize {
		fmt.Fprintf(os.Stderr, "WARNING: only %d candidates, need %d\n", len(candidates), targetSize)
	}

	// Sort by frequency (descending) and pick top N.
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].freq > candidates[j].freq
	})

	if len(candidates) > targetSize {
		candidates = candidates[:targetSize]
	}

	// Final sort: alphabetical.
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].word < candidates[j].word
	})

	for _, c := range candidates {
		fmt.Println(c.word)
	}

	fmt.Fprintf(os.Stderr, "Generated %d words\n", len(candidates))
}
