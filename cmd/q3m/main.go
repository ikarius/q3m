package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "q3m",
	Short: "q3m - géocodage en 3 mots pour la France métropolitaine",
	Long: `q3m encode les coordonnées GPS en triplets de mots français
sur une grille Lambert93 de précision 1m x 1m.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "sortie au format JSON")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
