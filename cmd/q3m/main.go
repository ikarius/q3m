package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "q3m",
	Short: "q3m - geocodage en 3 mots pour la France metropolitaine",
	Long: `q3m encode les coordonnees GPS en triplets de mots francais
sur une grille Lambert93 de precision 1m x 1m.`,
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
