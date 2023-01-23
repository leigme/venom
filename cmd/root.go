package cmd

/*
Copyright © 2023 NAME HERE <leigme@gmail.com>
*/
import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "venom",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.venom.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func InitWorkDir() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	appName := trimSuffix(executable)
	workDir := fmt.Sprint(".", filepath.Base(appName))
	if userHome, err := os.UserHomeDir(); err == nil {
		workDir = filepath.Join(userHome, workDir)
	}
	err = os.MkdirAll(workDir, os.ModePerm)
	if err != nil {
		log.Fatalf("mkdir work dir is err: %s\n", err)
	}
}

func trimSuffix(filename string) string {
	if strings.EqualFold(runtime.GOOS, "windows") {
		filename = strings.TrimSuffix(filename, ".exe")
	}
	return filename
}
