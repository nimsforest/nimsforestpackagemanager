package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nimsforestpm",
	Short: "NimsForest Package Manager - Simple Go-based tool manager",
	Long: `NimsForest Package Manager is a lightweight tool manager that installs and manages 
NimsForest tools via go get and go install. No complex dependencies, no configuration files—
just a simple wrapper around Go's native tooling.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}