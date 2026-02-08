package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ikarius/q3m"
	"github.com/spf13/cobra"
)

var fromlamCmd = &cobra.Command{
	Use:   "fromlam <E> <N>",
	Short: "Convertit des coordonnées Lambert93 en WGS84",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		E, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erreur: coordonnée E invalide: %v\n", err)
			os.Exit(1)
		}
		N, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erreur: coordonnée N invalide: %v\n", err)
			os.Exit(1)
		}

		lat, lon := q3m.FromLambert93(E, N)

		if jsonOutput {
			out := struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
				E   float64 `json:"e"`
				N   float64 `json:"n"`
			}{
				Lat: lat,
				Lon: lon,
				E:   E,
				N:   N,
			}
			writeJSON(out)
		} else {
			fmt.Printf("%.6f, %.6f\n", lat, lon)
		}
	},
}

func init() {
	rootCmd.AddCommand(fromlamCmd)
}
