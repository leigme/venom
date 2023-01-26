package cmd

import (
	"fmt"
	"github.com/leigme/venom/tool"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	AddCommand(&ThorCommand{
		version: "latest",
	})
}

type ThorCommand struct {
	version string
}

func (tc *ThorCommand) Execute() Exec {
	return func(cmd *cobra.Command, args []string) {
		if len(args) != 0 && !strings.EqualFold(args[0], "") {
			tc.version = args[0]
		}
		command := fmt.Sprintf("go install github.com/leigme/thor@%s", tc.version)
		v := tool.NewVmd()
		fmt.Println(v.Execute(command))
	}
}
