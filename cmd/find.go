package cmd

import (
	"fmt"
	"github.com/leigme/venom/tool"
	"github.com/spf13/cobra"
)

func init() {
	AddCommand(&FindCommand{}, CommandWithShort("find"))
}

type FindCommand struct {
}

func (fc *FindCommand) Execute() Exec {
	return func(cmd *cobra.Command, args []string) {
		v := tool.NewVmd()
		fmt.Println(v.Execute("ifconfig"))
	}
}
