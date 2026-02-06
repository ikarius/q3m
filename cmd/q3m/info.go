package main

import (
	"fmt"

	"github.com/ikarius/q3m"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Affiche les paramètres de la grille q3m",
	Run: func(cmd *cobra.Command, args []string) {
		if jsonOutput {
			fmt.Printf(`{"projection":"Lambert93/EPSG:2154","emin":%d,"emax":%d,"nmin":%d,"nmax":%d,"grid_width":%d,"grid_height":%d,"total_cells":%d,"dict_size":%d,"precision":"1m x 1m"}`+"\n",
				int(q3m.EMin), int(q3m.EMax), int(q3m.NMin), int(q3m.NMax),
				q3m.GridWidth, q3m.GridHeight, q3m.TotalCells, q3m.DictSize)
		} else {
			fmt.Println("q3m - géocodage en 3 mots")
			fmt.Println()
			fmt.Println("Projection:    Lambert93 / EPSG:2154 (GRS80)")
			fmt.Printf("Emprise E:     %d - %d m\n", int(q3m.EMin), int(q3m.EMax))
			fmt.Printf("Emprise N:     %d - %d m\n", int(q3m.NMin), int(q3m.NMax))
			fmt.Printf("Largeur:       %d cellules\n", q3m.GridWidth)
			fmt.Printf("Hauteur:       %d cellules\n", q3m.GridHeight)
			fmt.Printf("Total:         %d cellules\n", q3m.TotalCells)
			fmt.Printf("Dictionnaire:  %d mots\n", q3m.DictSize)
			fmt.Println("Précision:     1m x 1m")
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
