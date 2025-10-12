package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Run(cmd *cobra.Command, args []string) error {
	toggle, err := cmd.Flags().GetBool("toggle")
	if err != nil {
		return err
	}

	if toggle {
		fmt.Println("Hey, its Spot with toggle on!")
		return nil
	}

	fmt.Println("Hey, its Spot with toggle off!")
	return nil
}
