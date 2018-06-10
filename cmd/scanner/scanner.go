package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/bgo-education/test-grader-client/pkg/option"
	"github.com/bgo-education/test-grader-client/pkg/scan"
	"github.com/bgo-education/test-grader-client/pkg/utils"
)

var opt = option.GetInstance()

func main() {
	if !Auth() {
		fmt.Println("Wrong username or password!")
		return
	}
	fmt.Println("Authentication successful!")

	defer utils.Duration(time.Now(), "Scanner")

	fmt.Printf("Scanning %s...\n", opt.SrcDirectory)
	folders, err := ioutil.ReadDir(opt.SrcDirectory)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Found %d files\n", len(folders))

	count, fail := 0, 0
	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}

		name := folder.Name()

		if !opt.Override && scan.CheckFolder(name) {
			fmt.Printf("Skip %s", name)
			continue
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

		fmt.Printf("Processing %s\n", name)
		err := scan.ProcessFolder(name, writeChan)
		if err != nil {
			fmt.Println(err)
			fail++
			close(writeChan)
			continue
		}
		close(writeChan)
		done <- 1

		err = report.ToCSV(opt.Dst)
		if err != nil {
			fmt.Println(err)
			fail++
			continue
		}

		count++
	}

	fmt.Printf("Process success %d folders, failed %d\n", count, fail)
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
