package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// notifierCmd represents the notifier command
var notifierCmd = &cobra.Command{
	Use:   "notifier",
	Short: "It's responsible for notification purposes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("notifier called")
	},
}

func init() {
	rootCmd.AddCommand(notifierCmd)
}
