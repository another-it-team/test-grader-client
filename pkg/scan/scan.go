package scan

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bgo-education/test-grader-client/pkg/option"
	"github.com/bgo-education/test-grader-client/pkg/utils"
	"github.com/sirupsen/logrus"
)

const (
	POST = "POST"
	GET  = "GET"
)

var opt = option.GetInstance()
var logger = logrus.WithField("module", "scan")

func ProcessFolder(folder string, writeChan chan<- []string) error {
	files := utils.GetFilesByType(folder, opt.FilesExtension)
	logger.Infof("Found %d files", len(files))

	client := &http.Client{}

	for _, file := range files {
		if opt.Verbose {
			logger.Infof("Read %s", file)
		}

		req, err := UploadFile(file)
		if req != nil && err == nil {
			res, err := client.Do(req)
			if err != nil {
				logger.Error(err)
				continue
			}

			body := &bytes.Buffer{}
			_, err = body.ReadFrom(res.Body)
			if err != nil {
				logger.Error(err)
			}
			res.Body.Close()

			var data []string
			err = json.Unmarshal(body.Bytes(), &data)
			if err != nil {
				logger.Error(err)
			}

			writeChan <- data

			if opt.Verbose {
				logger.Infof("File %s, status code: %d", file, res.StatusCode)
			}
		} else {
			logger.Error(err)
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
