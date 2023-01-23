package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(renameCommand)
}

var renameCommand = &cobra.Command{Use: "rename"}
