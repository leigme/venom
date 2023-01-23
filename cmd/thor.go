package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(thorCommand)
}

var thorCommand = &cobra.Command{
	Use: "thor",
}
