package main

import (
	"fmt"
	"os"

	"github.com/bthompso/engx-ergonomics-poc/internal/commands"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "engx",
		Short: "ENGX - React Application Development Tool (POC)",
		Long: `ENGX POC - A terminal-based simulation of React web application creation.
This tool demonstrates human-computer interaction patterns for developer tooling.

Focus: Terminal UI/UX patterns, not actual application scaffolding.`,
		Version: fmt.Sprintf("%s (%s) built on %s", version, commit, date),
	}

	// Add global verbosity flags (mutually exclusive)
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Show only essential information")
	rootCmd.PersistentFlags().Bool("concise", false, "Show less detail with components hidden")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show enhanced details with progress bars for multi-step")
	rootCmd.PersistentFlags().Bool("debug", false, "Show maximum verbosity with all system outputs")

	// Other global flags
	rootCmd.PersistentFlags().String("config", "", "Config file (default searches for .engx/config.yaml)")

	// Mark verbosity flags as mutually exclusive
	rootCmd.MarkFlagsMutuallyExclusive("quiet", "concise", "verbose", "debug")

	// Add commands
	rootCmd.AddCommand(commands.NewCreateCommand())
	rootCmd.AddCommand(commands.NewTestErrorCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}