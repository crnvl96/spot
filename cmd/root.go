package cmd

import (
	"os"

	"github.com/crnvl96/spot/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spot",
	Short: "Spot is a cli tool to check git repositories for uncommitted or unpushed changes",
	Long: `
Spot scans specified directories for git repositories and checks for uncommitted changes or unpushed commits.

Use -t to specify target directories. Append /** to include subdirectories up to the depth specified by -d.
Examples:
  spot -t ~/config/nvim ~/Developer/**
  spot -d 2 -t ~/Developer/**
`,
	RunE: internal.Run,
}

func init() {
	rootCmd.Flags().IntP("depth", "d", 1, "Depth for subdirectory search when using /**")
	rootCmd.Flags().StringSliceP("target", "t", []string{}, "Target directories to scan")
	rootCmd.MarkFlagRequired("target")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
