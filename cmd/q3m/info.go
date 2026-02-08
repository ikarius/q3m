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
			out := struct {
				Projection string `json:"projection"`
				EMin       int    `json:"emin"`
				EMax       int    `json:"emax"`
				NMin       int    `json:"nmin"`
				NMax       int    `json:"nmax"`
				GridWidth  uint64 `json:"grid_width"`
				GridHeight uint64 `json:"grid_height"`
				TotalCells uint64 `json:"total_cells"`
				DictSize   int    `json:"dict_size"`
				Precision  string `json:"precision"`
			}{
				Projection: "Lambert93/EPSG:2154",
				EMin:       int(q3m.EMin),
				EMax:       int(q3m.EMax),
				NMin:       int(q3m.NMin),
				NMax:       int(q3m.NMax),
				GridWidth:  q3m.GridWidth,
				GridHeight: q3m.GridHeight,
				TotalCells: q3m.TotalCells,
				DictSize:   q3m.DictSize,
				Precision:  "1m x 1m",
			}
			writeJSON(out)
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
