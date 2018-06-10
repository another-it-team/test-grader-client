package main

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/kietdinh/test-grader-client/pkg/option"
	"github.com/kietdinh/test-grader-client/pkg/scan"
	"github.com/kietdinh/test-grader-client/pkg/utils"
	"github.com/sirupsen/logrus"
)

var opt = option.GetInstance()
var logger = logrus.WithField("module", "main")

func init() {
	err := os.MkdirAll(opt.DstDirectory, os.ModePerm)
	if err != nil {
		logger.Error(err)
	}
}

func main() {
	defer utils.Duration(time.Now(), "Scanner")

	logger.Info("Create session...")
	id, err := scan.CreateSession()
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("Session %d was created", id)

	logger.Infof("Scanning %s...", opt.SrcDirectory)
	folders, err := ioutil.ReadDir(opt.SrcDirectory)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Infof("Found %d files", len(folders))

	count, fail := 0, 0
	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}

		name := folder.Name()
		logger.Infof("Processing %s", name)
		err := scan.ProcessFolder(name, id)
		if err != nil {
			logger.Error(err)
			fail++
			continue
		}

		count++
	}

	logger.Infof("Process success %d folders, failed %d", count, fail)

	err = scan.CloseSession(id)
	if err != nil {
		logger.Error(err)
	}

	err = scan.GetResultFile(id)
	if err != nil {
		logger.Error(err)
	}
}
