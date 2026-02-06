# q3m - Géocodage en 3 mots pour la France métropolitaine

**q3m** encode n'importe quelle position GPS en France métropolitaine (Corse incluse) en un triplet de trois mots français, avec une précision de **1 mètre**.

```
48.8584, 2.2945  -->  province.shootons.retirons
                 <--  48.858398, 2.294503
```

## Pourquoi q3m ?

Le système what3words découpe le globe en cellules de 3m x 3m sur WGS84. En réalité, comme les degrés de longitude rétrécissent vers les pôles, ces cellules ne sont pas carrées.

q3m résout ce problème en utilisant la projection **Lambert93** (EPSG:2154), une projection métrique officielle de l'IGN. Chaque cellule mesure exactement **1m x 1m** dans le plan projeté.

## Installation

```bash
go install github.com/ikarius/q3m/cmd/q3m@latest
```

Ou depuis les sources :

```bash
git clone https://github.com/ikarius/q3m.git
cd q3m
go build ./cmd/q3m/
```

## Utilisation CLI

### Encoder des coordonnées

```bash
q3m encode 48.8584 2.2945
# province.shootons.retirons
```

### Décoder une adresse

```bash
q3m decode province.shootons.retirons
# 48.858398, 2.294503
```

### Informations de la grille

```bash
q3m info
```

### Sortie JSON

Toutes les commandes acceptent le flag `--json` :

```bash
q3m encode 48.8584 2.2945 --json
# {"address":"province.shootons.retirons","w1":"province","w2":"shootons","w3":"retirons","lat":48.858400,"lon":2.294500}

q3m decode province.shootons.retirons --json
# {"lat":48.858398,"lon":2.294503,"address":"province.shootons.retirons"}
```

## Utilisation comme bibliothèque Go

```go
import "github.com/ikarius/q3m"

// Encoder
addr, err := q3m.Encode(48.8584, 2.2945)
fmt.Println(addr) // province.shootons.retirons

// Décoder
coord, err := q3m.Decode("province.shootons.retirons")
fmt.Printf("%.6f, %.6f\n", coord.Lat, coord.Lon)
```

### API

| Fonction | Signature | Description |
|---|---|---|
| `Encode` | `(lat, lon float64) -> (Address, error)` | Coordonnées GPS vers adresse q3m |
| `Decode` | `(address string) -> (Coordinate, error)` | Adresse q3m vers coordonnées GPS |
| `ToLambert93` | `(lat, lon float64) -> (E, N float64)` | WGS84 vers Lambert93 |
| `FromLambert93` | `(E, N float64) -> (lat, lon float64)` | Lambert93 vers WGS84 |

### Types

```go
type Coordinate struct {
    Lat float64
    Lon float64
}

type Address struct {
    W1 string
    W2 string
    W3 string
}
```

## Paramètres techniques

| Paramètre | Valeur |
|---|---|
| Projection | Lambert93 / EPSG:2154 (ellipsoïde GRS80) |
| Emprise E | 100 000 - 1 250 000 m |
| Emprise N | 6 050 000 - 7 120 000 m |
| Grille | 1 150 000 x 1 070 000 cellules |
| Total | 1 230 500 000 000 cellules (~1.23 x 10^12) |
| Dictionnaire | 10 800 mots (10 800^3 = 1.26 x 10^12) |
| Précision | 1m x 1m (erreur max 0.71m du centre au coin) |
| Couverture | France métropolitaine + Corse |

## Comment ça marche

### Encodage

1. Les coordonnées WGS84 `(lat, lon)` sont projetées en Lambert93 `(E, N)`
2. La position est discrétisée en cellule de 1m x 1m : `x = floor(E - E_min)`, `y = floor(N - N_min)`
3. Un index linéaire est calculé : `idx = y * largeur + x`
4. L'index est permuté par un réseau de Feistel (décorrélation spatiale)
5. L'index permuté est converti en base 10 800 : trois indices de mots
6. Chaque indice est remplacé par le mot correspondant dans le dictionnaire

### Décodage

Le processus inverse exact. Le centre de la cellule (+0.5m) est retourné.

### Décorrélation spatiale

Sans la permutation, deux points voisins auraient des adresses presque identiques (deux mots sur trois en commun). Le réseau de Feistel assure que des cellules adjacentes produisent des triplets complètement différents, ce qui réduit les risques de confusion.

## Dictionnaire

Les 10 800 mots sont extraits de **Lexique383** (lexique.org), une base lexicale française libre.

Critères de sélection :
- 4 à 8 lettres
- Pas d'accents (ASCII uniquement)
- Noms, adjectifs, verbes, adverbes
- Triés par fréquence d'usage, les plus courants en priorité

Le dictionnaire est embarqué dans le binaire via `go:embed`. L'outil `tools/wordgen/` permet de régénérer le fichier `words_fr.txt` à partir de Lexique383.

**Contrat de stabilité** : une fois figé en v1.0, le dictionnaire et la clé de permutation ne changent plus jamais. Toute modification invaliderait les adresses existantes.

## Performance

Mesurée sur AMD Ryzen 9 8945HS :

| Opération | Temps | Allocations |
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

## Structure du projet

```
q3m/
├── go.mod                 # Module Go
├── lambert93.go           # Projection Lambert93 <-> WGS84
├── lambert93_test.go
├── grid.go                # Grille 1m, indexation cellules
├── grid_test.go
├── shuffle.go             # Permutation Feistel (décorrélation)
├── shuffle_test.go
├── words.go               # Dictionnaire (go:embed, sync.Once)
├── words_test.go
├── words_fr.txt           # 10 800 mots français
├── q3m.go                 # API publique : Encode(), Decode()
├── q3m_test.go
├── cmd/q3m/
│   ├── main.go            # Point d'entrée CLI (Cobra)
│   ├── encode.go          # Sous-commande encode
│   ├── decode.go          # Sous-commande decode
│   └── info.go            # Sous-commande info
└── tools/wordgen/
    └── main.go            # Génération du dictionnaire (Lexique383)
```

## Limitations

- **Couverture** : France métropolitaine et Corse uniquement. Les DOM-TOM ne sont pas couverts par Lambert93.
- **Cellules en mer** : tout le rectangle englobant Lambert93 est encodé, y compris les zones maritimes.
- **Pas de correction orthographique** : un mot mal saisi retournera une erreur, pas une suggestion.

## Licence

Ce projet est distribué sous licence [Mozilla Public License 2.0](LICENSE).

Le dictionnaire (`words_fr.txt`) est dérivé de Lexique383, distribué sous [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).

## Crédits

- **Lexique383** (lexique.org) pour la base lexicale
- **IGN** pour les paramètres de la projection Lambert93/RGF93
