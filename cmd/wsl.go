package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(wslCommand)
}

var wslCommand = &cobra.Command{Use: "wsl"}
