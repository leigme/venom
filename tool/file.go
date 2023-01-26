package tool

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var workDir string

func GetWorkDir() string {
	if strings.EqualFold(workDir, "") {
		appName := getAppName()
		workDir = fmt.Sprint(dirPrefix, filepath.Base(appName))
		if userHome, err := os.UserHomeDir(); err == nil {
			workDir = filepath.Join(userHome, workDir)
		}
	}
	return workDir
}

func getAppName() string {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	if strings.EqualFold(runtime.GOOS, windowsOs) {
		executable = strings.TrimSuffix(executable, exeSuffix)
	}
	return filepath.Base(executable)
}

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func CreateDir(filename string) error {
	if FileExist(filename) {
		return nil
	}
	return os.MkdirAll(filename, os.ModePerm)
}

func FileMD5(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func Bytes2File(bytes []byte, filename string) {
	var (
		f   *os.File
		err error
	)
	if !FileExist(filename) {
		f, err = os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		f, err = os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer f.Close()
	fw := bufio.NewWriter(f)
	_, err = fw.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
	err = fw.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

type FileList []fs.DirEntry

func (fl FileList) Len() int {
	return len(fl)
}

func (fl FileList) Less(i, j int) bool {
	fsi, _ := fl[i].Info()
	fsj, _ := fl[j].Info()
	return fsi.ModTime().After(fsj.ModTime())
}

func (fl FileList) Swap(i, j int) {
	fl[i], fl[j] = fl[j], fl[i]
}
