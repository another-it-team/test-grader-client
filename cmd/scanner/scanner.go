package main

import (
	"io/ioutil"
	"time"

	"github.com/bgo-education/test-grader-client/pkg/option"
	"github.com/bgo-education/test-grader-client/pkg/scan"
	"github.com/bgo-education/test-grader-client/pkg/utils"
	"github.com/sirupsen/logrus"
)

var opt = option.GetInstance()
var logger = logrus.WithField("module", "main")

func main() {
	defer utils.Duration(time.Now(), "Scanner")

	logger.Infof("Scanning %s...", opt.SrcDirectory)
	folders, err := ioutil.ReadDir(opt.SrcDirectory)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("Found %d files", len(folders))

	report := scan.NewReport(scan.Header(opt.NumCau))

	writeChan := make(chan []string, 1000)
	go func() {
		for d := range writeChan {
			report.Add(d)
		}
	}()

	count, fail := 0, 0
	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}

		name := folder.Name()
		logger.Infof("Processing %s", name)
		err := scan.ProcessFolder(name, writeChan)
		if err != nil {
			logger.Error(err)
			fail++
			continue
		}

		count++
	}

	logger.Infof("Process success %d folders, failed %d", count, fail)

	err = report.ToCSV(opt.Dst)
	if err != nil {
		logger.Error(err)
	}
}
