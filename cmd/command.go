package cmd

import (
	"github.com/leigme/loki/app"
	"github.com/leigme/loki/file"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

const goSuffix = ".go"

type CommandInterface interface {
	Execute() Exec
}

type Exec func(cmd *cobra.Command, args []string)

func AddCommand(ci CommandInterface, ops ...Option) {
	rootCmd.AddCommand(NewCommand(ci.Execute(), ops...))
}

type CommandOption struct {
	Short, Long string
}

type Option func(option *CommandOption)

func NewCommand(e Exec, ops ...Option) *cobra.Command {
	cmdName := getCommandName(3)
	err := file.CreateDir(filepath.Join(app.WorkDir(), cmdName))
	if err != nil {
		log.Fatal(err)
	}
	co := newDefaultOption()
	for _, apply := range ops {
		apply(co)
	}
	return &cobra.Command{Use: cmdName, Short: co.Short, Long: co.Long, Run: e}
}

func CommandWithShort(short string) Option {
	return func(option *CommandOption) {
		if strings.EqualFold(short, "") {
			return
		}
		option.Short = short
	}
}

func CommandWithLong(long string) Option {
	return func(option *CommandOption) {
		if strings.EqualFold(long, "") {
			return
		}
		option.Long = long
	}
}

func newDefaultOption() *CommandOption {
	return &CommandOption{
		Short: "",
		Long:  "",
	}
}

func getCommandName(skip int) string {
	_, filename, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	return strings.TrimSuffix(filepath.Base(filename), goSuffix)
}
