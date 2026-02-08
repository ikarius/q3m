package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags.
var version = "dev"

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:     "q3m",
	Short:   "q3m - géocodage en 3 mots pour la France métropolitaine",
	Long:    "q3m encode les coordonnées GPS en triplets de mots français\nsur une grille Lambert93 de précision 1m x 1m.",
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "sortie au format JSON")
	rootCmd.SilenceErrors = true
}

func writeJSON(v any) {
	if err := json.NewEncoder(os.Stdout).Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "erreur JSON: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
