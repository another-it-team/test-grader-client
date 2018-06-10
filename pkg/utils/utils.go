package utils

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("module", "duration")

func Duration(start time.Time, job string) {
	logger.Infof("Job %s cost %v", job, time.Now().Sub(start))
}

func GetFilesByType(folder string, exts []string) (files []string) {
	for _, dtype := range exts {
		fs, err := filepath.Glob(filepath.Join(folder, "*"+dtype))
		if err != nil {
			logger.Error(err)
			continue
		}

		files = append(files, fs...)
	}
	return
}

func DownloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
