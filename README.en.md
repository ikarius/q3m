# q3m - 3-word geocoding for metropolitan France

[![CI](https://github.com/ikarius/q3m/actions/workflows/ci.yml/badge.svg)](https://github.com/ikarius/q3m/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/ikarius/q3m)](https://github.com/ikarius/q3m/releases/latest)
[![License: MPL-2.0](https://img.shields.io/github/license/ikarius/q3m)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/ikarius/q3m.svg)](https://pkg.go.dev/github.com/ikarius/q3m)
[![Go Report Card](https://goreportcard.com/badge/github.com/ikarius/q3m)](https://goreportcard.com/report/github.com/ikarius/q3m)

**q3m** encodes any GPS position in metropolitan France (including Corsica) into a triplet of three French words, with a precision of **1 metre**.

```
48.8584, 2.2945  -->  province.shootons.retirons
                 <--  48.858398, 2.294503
```

## Why q3m?

The what3words system divides the globe into 3m x 3m cells on WGS84. In practice, since longitude degrees shrink towards the poles, these cells are not actually square.

q3m solves this by using the **Lambert93** projection (EPSG:2154), the official metric projection from the French National Geographic Institute (IGN). Each cell measures exactly **1m x 1m** in the projected plane.

## Installation

```bash
go install github.com/ikarius/q3m/cmd/q3m@latest
```

Or from source:

```bash
git clone https://github.com/ikarius/q3m.git
cd q3m
go build ./cmd/q3m/
```

## CLI usage

### Encode coordinates

```bash
q3m encode 48.8584 2.2945
# province.shootons.retirons
```

### Decode an address

```bash
q3m decode province.shootons.retirons
# 48.858398, 2.294503
```

### Grid information

```bash
q3m info
```

### JSON output

All commands accept the `--json` flag:

```bash
q3m encode 48.8584 2.2945 --json
# {"address":"province.shootons.retirons","w1":"province","w2":"shootons","w3":"retirons","lat":48.858400,"lon":2.294500}

q3m decode province.shootons.retirons --json
# {"lat":48.858398,"lon":2.294503,"address":"province.shootons.retirons"}
```

## Go library usage

```go
import "github.com/ikarius/q3m"

// Encode
addr, err := q3m.Encode(48.8584, 2.2945)
fmt.Println(addr) // province.shootons.retirons

// Decode
coord, err := q3m.Decode("province.shootons.retirons")
fmt.Printf("%.6f, %.6f\n", coord.Lat, coord.Lon)
```

### API

| Function | Signature | Description |
|---|---|---|
| `Encode` | `(lat, lon float64) -> (Address, error)` | GPS coordinates to q3m address |
| `Decode` | `(address string) -> (Coordinate, error)` | q3m address to GPS coordinates |
| `ToLambert93` | `(lat, lon float64) -> (E, N float64)` | WGS84 to Lambert93 |
| `FromLambert93` | `(E, N float64) -> (lat, lon float64)` | Lambert93 to WGS84 |

### Types

```go
type Coordinate struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}

type Address struct {
    W1 string `json:"w1"`
    W2 string `json:"w2"`
    W3 string `json:"w3"`
}
```

## Technical parameters

| Parameter | Value |
|---|---|
| Projection | Lambert93 / EPSG:2154 (GRS80 ellipsoid) |
| Easting range | 100,000 - 1,250,000 m |
| Northing range | 6,050,000 - 7,120,000 m |
| Grid | 1,150,000 x 1,070,000 cells |
| Total | 1,230,500,000,000 cells (~1.23 x 10^12) |
| Dictionary | 10,800 words (10,800^3 = 1.26 x 10^12) |
| Precision | 1m x 1m (max error 0.71m from centre to corner) |
| Coverage | Metropolitan France + Corsica |

## How it works

### Encoding

1. WGS84 coordinates `(lat, lon)` are projected to Lambert93 `(E, N)`
2. The position is discretised into a 1m x 1m cell: `x = floor(E - E_min)`, `y = floor(N - N_min)`
3. A linear index is computed: `idx = y * width + x`
4. The index is permuted through a Feistel network (spatial decorrelation)
5. The permuted index is converted to base 10,800: three word indices
6. Each index is replaced by the corresponding word from the dictionary

### Decoding

The exact reverse process. The cell centre (+0.5m) is returned.

### Spatial decorrelation

Without the permutation, two neighbouring points would share nearly identical addresses (two words out of three in common). The Feistel network ensures that adjacent cells produce completely different triplets, reducing the risk of confusion.

## Dictionary

The 10,800 words are sourced from **Lexique383** (lexique.org), an open French lexical database.

Selection criteria:
- 4 to 8 letters
- No diacritical marks (ASCII only)
- Nouns, adjectives, verbs, adverbs
- Sorted by usage frequency, most common words first

The dictionary is embedded in the binary via `go:embed`. The `tools/wordgen/` tool can regenerate `words_fr.txt` from Lexique383.

**Stability contract**: once frozen at v1.0, the dictionary and permutation key never change. Any modification would invalidate all existing addresses.

## Performance

Measured on AMD Ryzen 9 8945HS:

| Operation | Time | Allocations |
|---|---|---|
| Encode | 133 ns/op | 0 |
| Decode | 665 ns/op | 1 |
| ToLambert93 | 71 ns/op | 0 |
| FromLambert93 | 482 ns/op | 0 |
| Shuffle | 103 ns/op | 0 |

## Tests

```bash
go test ./...
go test -bench . -benchmem
```

## Project structure

```
q3m/
├── go.mod                 # Go module
├── lambert93.go           # Lambert93 <-> WGS84 projection
├── lambert93_test.go
├── grid.go                # 1m grid, cell indexation
├── grid_test.go
├── shuffle.go             # Feistel permutation (decorrelation)
├── shuffle_test.go
├── words.go               # Dictionary (go:embed, sync.Once)
├── words_test.go
├── words_fr.txt           # 10,800 French words
├── q3m.go                 # Public API: Encode(), Decode()
├── q3m_test.go
├── cmd/q3m/
│   ├── main.go            # CLI entry point (Cobra)
│   ├── encode.go          # encode subcommand
│   ├── decode.go          # decode subcommand
│   └── info.go            # info subcommand
└── tools/wordgen/
    └── main.go            # Dictionary generation (Lexique383)
```

## Limitations

- **Coverage**: metropolitan France and Corsica only. Overseas territories are not covered by Lambert93.
- **Cells at sea**: the entire Lambert93 bounding rectangle is encoded, including maritime areas.
- **No spell-checking**: a misspelled word will return an error, not a suggestion.

## Licence

This project is licensed under the [Mozilla Public License 2.0](LICENSE).

The dictionary (`words_fr.txt`) is derived from Lexique383, distributed under [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).

## Credits

- **Lexique383** (lexique.org) for the lexical database
- **IGN** (French National Geographic Institute) for Lambert93/RGF93 projection parameters
