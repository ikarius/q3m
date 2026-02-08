package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ikarius/q3m"
	"github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
	Use:   "decode <mot1.mot2.mot3>",
	Short: "Décode une adresse q3m en coordonnées GPS",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		coord, err := q3m.Decode(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "erreur: %v\n", err)
			os.Exit(1)
		}

		if jsonOutput {
			addr := strings.ToLower(strings.TrimSpace(args[0]))
			parts := strings.Split(addr, ".")
			out := struct {
				Lat     float64 `json:"lat"`
				Lon     float64 `json:"lon"`
				Address string  `json:"address"`
				W1      string  `json:"w1"`
				W2      string  `json:"w2"`
				W3      string  `json:"w3"`
			}{
				Lat:     coord.Lat,
				Lon:     coord.Lon,
				Address: addr,
				W1:      parts[0],
				W2:      parts[1],
				W3:      parts[2],
			}
			writeJSON(out)
		} else {
			fmt.Printf("%.6f, %.6f\n", coord.Lat, coord.Lon)
		}
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)
}
