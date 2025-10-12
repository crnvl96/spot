package cmd

import (
	"os"

	"github.com/crnvl96/spot/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spot",
	Short: "Spot is a cli tool to check if you have uncommited changes in any of your repos",
	Long: `
Spot is a cli tool that scans all folders present in the current directory.

If any of these folders is a git repository, it will automatically check if it has any uncommited or unpushed changed, notifying you at the end.
The main goal of this tool it to help you prevent losing work by forgetting to commit changes to any of your repositories.
`,
	RunE: internal.Run,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
