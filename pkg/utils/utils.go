package utils

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				return err
			}
		}

		rc.Close()
	}
	return nil
}

func GetCurrentDir() (dir string, err error) {
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}
	temp := strings.Split(dir, "\\")
	dir = strings.ToLower(temp[len(temp)-1])
	return
}
