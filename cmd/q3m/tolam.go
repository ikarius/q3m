package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ikarius/q3m"
	"github.com/spf13/cobra"
)

var tolamCmd = &cobra.Command{
	Use:   "tolam <lat> <lon>",
	Short: "Convertit des coordonnées WGS84 en Lambert93",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		lat, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erreur: latitude invalide: %v\n", err)
			os.Exit(1)
		}
		lon, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erreur: longitude invalide: %v\n", err)
			os.Exit(1)
		}

		E, N := q3m.ToLambert93(lat, lon)

		if jsonOutput {
			out := struct {
				E   float64 `json:"e"`
				N   float64 `json:"n"`
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			}{
				E:   E,
				N:   N,
				Lat: lat,
				Lon: lon,
			}
			writeJSON(out)
		} else {
			fmt.Printf("%.4f, %.4f\n", E, N)
		}
	},
}

func init() {
	rootCmd.AddCommand(tolamCmd)
}
