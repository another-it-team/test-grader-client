package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Duration(start time.Time, job string) {
	fmt.Printf("Job %s cost %v\n", job, time.Now().Sub(start))
}

func GetFilesByType(folder string, exts []string) (files []string) {
	for _, dtype := range exts {
		fs, err := filepath.Glob(filepath.Join(folder, "*"+dtype))
		if err != nil {
			fmt.Println(err)
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

func ToMD5(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}
