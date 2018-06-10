package scan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bgo-education/test-grader-client/pkg/option"
	"github.com/bgo-education/test-grader-client/pkg/utils"
)

const (
	POST = "POST"
	GET  = "GET"
)

var opt = option.GetInstance()

func CheckFolder(folder string) bool {
	f, err := os.Open(filepath.Join(folder, opt.Dst))
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func ProcessFolder(folder string, writeChan chan<- []string) error {
	files := utils.GetFilesByType(folder, opt.FilesExtension)
	fmt.Printf("Found %d files\n", len(files))

	client := &http.Client{}

	for _, file := range files {
		if opt.Verbose {
			fmt.Printf("Read %s\n", file)
		}

		req, err := UploadFile(file)
		if req != nil && err == nil {
			res, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				continue
			}

			body := &bytes.Buffer{}
			_, err = body.ReadFrom(res.Body)
			if err != nil {
				fmt.Println(err)
				continue
			}
			res.Body.Close()

			var data []string
			err = json.Unmarshal(body.Bytes(), &data)
			if err != nil {
				fmt.Println(err)
				continue
			}

			writeChan <- data

			if opt.Verbose {
				fmt.Printf("File %s, status code: %d\n", file, res.StatusCode)
			}
		} else {
			fmt.Println(err)
		}
	}

	return nil
}

func UploadFile(filename string) (*http.Request, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, f)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(POST, opt.UploadEndPoint, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
