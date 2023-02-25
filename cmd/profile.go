package cmd

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/leigme/loki/app"
	loki "github.com/leigme/loki/cobra"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type ProfileCommand string

const (
	backup    ProfileCommand = "backup"
	restore   ProfileCommand = "restore"
	backupIni                = "backup.ini"
)

func init() {
	loki.Add(rootCmd, &profile{})
}

type profile struct {
	serverUrl string
}

func (p *profile) Execute() loki.Exec {
	return func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal("args must be at least 2")
		}
		switch ProfileCommand(args[0]) {
		case backup:
			cfg, err := ini.Load(filepath.Join(app.WorkDir(), "profile", backupIni))
			if err != nil {
				log.Fatal(err)
			}
			profilePath := cfg.Section("bash").Key("profile").String()
			aliases := cfg.Section("bash").Key("aliases").String()
			dir := cfg.Section("shell").Key("dir").String()
			uc := &UploadCommand{HttpClient: &http.Client{Timeout: 30 * time.Second}}
			if uc.CheckServer(fmt.Sprintf(FormatRunning, args[1])) {
				if !strings.EqualFold(profilePath, "") {
					uc.UploadFile(profilePath, fmt.Sprintf(FormatUpload, args[1]))
				}
				if !strings.EqualFold(aliases, "") {
					uc.UploadFile(aliases, fmt.Sprintf(FormatUpload, args[1]))
				}
				if !strings.EqualFold(dir, "") {
					uc.UploadFile(dir, fmt.Sprintf(FormatUpload, args[1]))
				}
			}

		case restore:

		}
	}
}
