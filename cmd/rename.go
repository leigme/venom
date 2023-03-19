package cmd

import (
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/leigme/loki/app"
	loki "github.com/leigme/loki/cobra"
	"github.com/leigme/loki/file"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type renameParam int

const (
	fileDir renameParam = iota
	regular
	suffixes
	section
)

const (
	renameCacheDir = "rename"
	regularConfig  = "regular.ini"
	suffixSplit    = "|"
)

func init() {
	loki.Add(rootCmd, &RenameCommand{params: make(map[renameParam]string, 0)})
}

type RenameCommand struct {
	params map[renameParam]string
}

func (rc *RenameCommand) Execute() loki.Exec {
	return func(cmd *cobra.Command, args []string) {
		rc.params[fileDir] = regularConfig
		rc.params[suffixes] = "*"
		rc.params[section] = "default"
		for i, v := range args {
			rc.params[renameParam(i)] = v
		}
		if strings.EqualFold(rc.params[fileDir], regularConfig) {
			showCacheRegular()
			return
		}
		if !file.Exist(rc.params[fileDir]) {
			log.Fatal("file dir is no exist")
		}
		if strings.EqualFold(rc.params[regular], "") {
			r, err := getCacheRegular(rc.params[section], rc.params[fileDir])
			if err != nil || strings.EqualFold(r, "") {
				log.Fatal(err)
			}
			rc.params[regular] = r
		} else {
			go setCacheRegular(rc.params[section], rc.params[fileDir], rc.params[regular])
		}
		reg, errReg := regexp.Compile(rc.params[regular])
		if errReg != nil {
			log.Fatal(errReg)
		}
		rc.changeFileName(reg)
	}
}

/*
设置缓存正则
*/
func setCacheRegular(section, key, value string) {
	filename := filepath.Join(app.WorkDir(), renameCacheDir, regularConfig)
	cfg := ini.Empty()
	defaultSection := cfg.Section(section)
	_, err := defaultSection.NewKey(key, value)
	if err != nil {
		log.Println(err)
	}
	err = cfg.SaveTo(filename)
	if err != nil {
		log.Println(err)
	}
}

/*
获取缓存正则
*/
func getCacheRegular(section, fileDir string) (string, error) {
	filename := filepath.Join(app.WorkDir(), renameCacheDir, regularConfig)
	if !file.Exist(filename) {
		return "", errors.New(fmt.Sprintf("%s is not exist\n", filename))
	}
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	return cfg.Section(section).Key(fileDir).Value(), nil
}

/*
显示缓存正则
*/
func showCacheRegular() {
	filename := filepath.Join(app.WorkDir(), renameCacheDir, regularConfig)
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	cfg.WriteTo(os.Stdout)
}

/*
根据文件后缀过滤
*/
func hasFileSuffix(path, suffixes string) bool {
	if strings.EqualFold(suffixes, "*") {
		return true
	}
	array := strings.Split(suffixes, suffixSplit)
	for _, s := range array {
		if strings.HasSuffix(path, s) {
			return true
		}
	}
	return false
}

/*
修改文件名
*/
func (rc *RenameCommand) changeFileName(reg *regexp.Regexp) {
	rc.params[fileDir] = strings.TrimSpace(rc.params[fileDir])
	fis, err := os.ReadDir(rc.params[fileDir])
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, fi := range fis {
		fn := fi.Name()
		fnSuffix := path.Ext(fn)
		if hasFileSuffix(fnSuffix, rc.params[suffixes]) {
			nfn := strings.ToUpper(reg.FindString(fn)) + path.Ext(fn)
			fna := filepath.Join(rc.params[fileDir], fi.Name())
			err := os.Rename(fna, filepath.Join(rc.params[fileDir], nfn))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
