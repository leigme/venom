package cmd

import (
	"fmt"
	loki "github.com/leigme/loki/cobra"
	"github.com/leigme/loki/shell"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	loki.Add(rootCmd, &ThorCommand{
		version: "latest",
	})
}

type ThorCommand struct {
	version string
}

func (tc *ThorCommand) Execute() loki.Exec {
	return func(cmd *cobra.Command, args []string) {
		if len(args) != 0 && !strings.EqualFold(args[0], "") {
			tc.version = args[0]
		}
		command := fmt.Sprintf("go install github.com/leigme/thor@%s", tc.version)
		sh := shell.New()
		fmt.Println(sh.Exe(command))
	}
}
