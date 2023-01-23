package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(findCommand)
}

var findCommand = &cobra.Command{Use: "find"}
