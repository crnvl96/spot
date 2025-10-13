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
Spot scans specified directories (and their subdirectories up to 2 levels deep) for git repositories and checks for uncommitted changes or unpushed commits.

If no targets are specified, it scans the current directory.

Examples:
  spot -t ~/config ~/Developer
`,
	RunE: internal.Run,
}

func init() {
	rootCmd.Flags().StringSliceP("target", "t", []string{}, "Target directories to scan (scans recursively up to 2 levels)")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
