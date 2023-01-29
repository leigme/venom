package cmd

import (
	"bytes"
	"fmt"
	"github.com/leigme/loki/file"
	"github.com/spf13/cobra"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	FormatRunning = "%s/running"
	FormatUpload  = "%s/upload"
)

func init() {
	AddCommand(&UploadCommand{
		HttpClient: &http.Client{Timeout: 10 * time.Minute},
	})
}

type UploadCommand struct {
	HttpClient *http.Client
}

func (uc *UploadCommand) Execute() Exec {
	return func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal("filename and url is not nil")
		}
		if uc.CheckServer(fmt.Sprintf(FormatRunning, args[1])) {
			uc.UploadFile(args[0], fmt.Sprintf(FormatUpload, args[1]))
		}
	}
}

func (uc *UploadCommand) CheckServer(url string) bool {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return false
	}
	resp, err := uc.HttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	return resp.StatusCode == http.StatusOK
}

func (uc *UploadCommand) UploadFile(filename, url string) {
	fi, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	srcMD5, err := file.Md5(filename)
	if err != nil {
		log.Fatalf("sum file: %s, md5 err: %s\n", filename, err)
	}
	log.Println(srcMD5)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		log.Fatal(err)
	}

	srcFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}

	err = writer.WriteField("md5", srcMD5)
	if err != nil {
		log.Fatal(err)
	}
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := uc.HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(content))
}
