package cmd

import (
	"errors"
	"fmt"
	"github.com/leigme/loki/app"
	loki "github.com/leigme/loki/cobra"
	"github.com/leigme/loki/file"
	"github.com/leigme/loki/shell"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type wslParam int

const (
	wslCmd wslParam = iota
	wslName
	wslFile
	wslDir
)

func init() {
	loki.Add(rootCmd, &WslCommand{params: make(map[wslParam]string, 0)})
}

const (
	wslWorkDir = "wsl"
	wslBackup  = "backup"
)

type WslCommand struct {
	params map[wslParam]string
}

func (wc *WslCommand) Execute() loki.Exec {
	return func(cmd *cobra.Command, args []string) {
		for i, v := range args {
			wc.params[wslParam(i)] = v
		}
		command := generateWslCmd(wc.params)
		if strings.EqualFold(command, "") {
			fmt.Println("command is nil")
			return
		}
		fmt.Println(command)
		sh := shell.New()
		fmt.Println(sh.Exe(command))
	}
}

func createDefaultExportFile(linux string) string {
	filename := fmt.Sprintf("%s_%s.tar", linux, time.Now().Format("200601021504"))
	if err := file.CreateDir(filepath.Join(app.WorkDir(), wslWorkDir, wslBackup)); err == nil {
		filename = filepath.Join(app.WorkDir(), wslWorkDir, wslBackup, filename)
	}
	return filename
}

func searchDefaultExportFile(linux string) (string, error) {
	path := filepath.Join(app.WorkDir(), wslWorkDir, wslBackup)
	if !file.Exist(path) {
		err := errors.New(fmt.Sprintf("%s path is not exist", path))
		return "", err
	}
	var fl file.ListWithTime
	fs, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	for _, f := range fs {
		if strings.HasPrefix(f.Name(), linux) && strings.HasSuffix(f.Name(), ".tar") {
			fl = append(fl, f)
		}
	}
	sort.Sort(fl)
	filename := filepath.Join(path, fl[0].Name())
	return filename, nil
}

func generateWslCmd(params map[wslParam]string) string {
	if strings.EqualFold(params[wslName], "") {
		log.Fatal("linux name is nil")
	}
	switch params[wslCmd] {
	case "export":
		if strings.EqualFold(params[wslFile], "") {
			params[wslFile] = createDefaultExportFile(params[wslName])
		}
		return generateExportCmd(params[wslName], params[wslFile])
	case "import":
		if strings.EqualFold(params[wslDir], "") {
			params[wslDir] = filepath.Join(app.WorkDir(), wslWorkDir)
		}
		if strings.EqualFold(params[wslFile], "") {
			fp, err := searchDefaultExportFile(params[wslName])
			if err != nil {
				log.Fatal(err)
			}
			params[wslFile] = fp
		}
		return generateImportCmd(params[wslName], params[wslDir], params[wslFile])
	default:
		return ""
	}
}

// wsl --export [Ubuntu-22.04] [C:\Users\leig\.wsl\backup\Ubuntu-22.04_20060102150405.tar]
func generateExportCmd(linux, exportFile string) string {
	return fmt.Sprintf("wsl --export %s %s", linux, exportFile)
}

// wsl --import [Ubuntu-22.04] [C:\Users\leig\.wsl\] [C:\Users\leig\.wsl\backup\Ubuntu-22.04_20060102150405.tar
func generateImportCmd(linux, importDir, importFile string) string {
	return fmt.Sprintf("wsl --import %s %s %s", linux, importDir, importFile)
}
