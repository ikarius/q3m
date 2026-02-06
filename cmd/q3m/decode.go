package main

import (
	"encoding/json"
	"fmt"
	"os"

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
			out := struct {
				Lat     float64 `json:"lat"`
				Lon     float64 `json:"lon"`
				Address string  `json:"address"`
			}{
				Lat:     coord.Lat,
				Lon:     coord.Lon,
				Address: args[0],
			}
			json.NewEncoder(os.Stdout).Encode(out)
		} else {
			fmt.Printf("%.6f, %.6f\n", coord.Lat, coord.Lon)
		}
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)
}
