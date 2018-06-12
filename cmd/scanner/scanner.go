package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bgo-education/test-grader-client/pkg/option"
	"github.com/bgo-education/test-grader-client/pkg/scan"
	"github.com/bgo-education/test-grader-client/pkg/utils"
)

var opt = option.GetInstance()

func main() {
	for !Auth() {
		fmt.Println("Wrong username or password!")
	}
	fmt.Println("Authentication successful!")

	defer utils.Duration(time.Now(), "Scanner")

	fmt.Println("Create session...")
	id, err := scan.CreateSession()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Session %s was created\n", id)

	fmt.Printf("Scanning %s...\n", opt.SrcDirectory)
	folders, err := ioutil.ReadDir(opt.SrcDirectory)
	if err != nil {
		fmt.Println(err)
		return
	}

	temp := strings.Split(opt.Dst, ".")
	dtype := temp[len(temp)-1]

	wg := &sync.WaitGroup{}
	count := 0
	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}

		wg.Add(1)
		go func(src string) {
			if !opt.Override && scan.CheckFolder(src) {
				fmt.Printf("Skip %s", src)
				return
			}

			report := scan.NewReport(scan.Header(opt.NumCau))
			writeChan := make(chan []string, 50)
			done := make(chan int)
			go func() {
				for d := range writeChan {
					report.Add(d)
				}
				<-done
			}()

			fmt.Printf("Processing %s\n", src)
			err := scan.ProcessFolder(src, id, writeChan)
			if err != nil {
				fmt.Println(err)
				close(writeChan)
				return
			}
			close(writeChan)
			done <- 1

			if dtype == scan.CSV {
				err = report.ToCSV(filepath.Join(src, opt.Dst))
			} else {
				err = report.ToXLSX(filepath.Join(src, opt.Dst))
			}
			if err != nil {
				fmt.Println(err)
				return
			}

			wg.Done()
		}(folder.Name())

		count++
	}

	wg.Wait()
	fmt.Printf("Process success %d folders\n", count)

	fmt.Println("Getting zip result file...")
	err = scan.GetImagesResult(scan.ImagesResult, id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Scanln() // wait for Enter Key
}

func Auth() bool {
	sc := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter username: ")
	sc.Scan()
	username := sc.Text()

	fmt.Print("Enter password: ")
	sc.Scan()
	password := utils.ToMD5(sc.Text())

	if username == opt.Username && password == opt.Password {
		return true
	}

	return false
}
