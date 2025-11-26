package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	noColor bool
	quiet   bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "env-diff [file1] [file2]",
		Short: "Compare .env files and detect configuration differences",
		Long:  "Compare two .env files to identify missing, extra, or mismatched variables. Useful for pre-deployment checks and keeping env templates in sync.",
		Args:  cobra.ExactArgs(2),
		Run:   runDiff,
	}

	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Only show summary, suppress detailed diff")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runDiff(cmd *cobra.Command, args []string) {
	file1, file2 := args[0], args[1]

	env1, err := ParseEnvFile(file1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", file1, err)
		os.Exit(1)
	}

	env2, err := ParseEnvFile(file2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", file2, err)
		os.Exit(1)
	}

	diff := CompareEnvFiles(env1, env2, file1, file2)
	diff.Print(noColor, quiet)

	if diff.HasDifferences() {
		os.Exit(1)
	}
}
