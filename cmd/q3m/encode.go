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
	Short: "Encode des coordonnees GPS en adresse q3m",
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
			fmt.Printf(`{"address":"%s","w1":"%s","w2":"%s","w3":"%s","lat":%f,"lon":%f}`+"\n",
				addr, addr.W1, addr.W2, addr.W3, lat, lon)
		} else {
			fmt.Println(addr)
		}
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)
}
