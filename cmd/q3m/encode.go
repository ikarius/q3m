package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ikarius/q3m"
	"github.com/spf13/cobra"
)

var encodeCmd = &cobra.Command{
	Use:   "encode <lat> <lon>",
	Short: "Encode des coordonnées GPS en adresse q3m",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		lat, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "latitude invalide: %v\n", err)
			os.Exit(1)
		}
		lon, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "longitude invalide: %v\n", err)
			os.Exit(1)
		}

		addr, err := q3m.Encode(lat, lon)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erreur: %v\n", err)
			os.Exit(1)
		}

		if jsonOutput {
			out := struct {
				Address string  `json:"address"`
				W1      string  `json:"w1"`
				W2      string  `json:"w2"`
				W3      string  `json:"w3"`
				Lat     float64 `json:"lat"`
				Lon     float64 `json:"lon"`
			}{
				Address: addr.String(),
				W1:      addr.W1,
				W2:      addr.W2,
				W3:      addr.W3,
				Lat:     lat,
				Lon:     lon,
			}
			writeJSON(out)
		} else {
			fmt.Println(addr)
		}
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)
}
