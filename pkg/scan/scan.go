package scan

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

	RESULT = "result.xlsx"
)

var opt = option.GetInstance()
var logger = logrus.WithField("module", "scan")

func ProcessFolder(folder string, id int) error {
	files := utils.GetFilesByType(folder, opt.FilesExtension)

	client := &http.Client{}
	url := fmt.Sprintf("%s/%d", opt.UploadEndPoint, id)
	params := map[string]string{
		"folder": folder,
		"name":   "",
	}

	for _, file := range files {
		params["name"] = file

		if opt.Verbose {
			logger.Infof("Read %s", file)
		}

		req, err := UploadFile(file, params, url)
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

			if opt.Verbose {
				logger.Infof("File %s, status code: %d", file, res.StatusCode)
			}
		} else {
			logger.Error(err)
		}
	}

	return nil
}

func UploadFile(filename string, params map[string]string, url string) (*http.Request, error) {
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

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(POST, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func CreateSession() (id int, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(POST, opt.CreateSessionEndPoint, bytes.NewBuffer(nil))
	if err != nil {
		return
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body.Bytes(), &id)
	if err != nil {
		return
	}

	return
}

func CloseSession(id int) error {
	client := &http.Client{}

	url := fmt.Sprintf("%s/%d", opt.CloseSessionEndPoint, id)
	req, err := http.NewRequest(POST, url, bytes.NewBuffer(nil))
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("close session not success, got %d", res.StatusCode)
		return errors.New(msg)
	}

	return nil
}

func GetResultFile(id int) error {
	url := fmt.Sprintf("%s/%d", opt.GetResultEndPoint, id)

	err := utils.DownloadFile(RESULT, url)
	if err != nil {
		return err
	}

	return nil
}
