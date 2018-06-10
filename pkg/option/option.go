package option

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bgo-education/test-grader-client/pkg/utils"
)

const (
	JPG  = ".jpg"
	PNG  = ".png"
	JPEG = ".jpeg"

	config    = "option.json"
	configURL = "https://bgo.edu.vn/test-grader-config.json"
)

var option *Option

type Option struct {
	// Auth
	Username string
	Password string

	Domain string

	// EndPoint API
	UploadEndPoint        string
	CreateSessionEndPoint string
	DownloadEndPoint      string

	// Dir
	SrcDirectory string
	Dst          string

	FilesExtension []string

	// SL Questions
	NumCau int

	Override bool
	Verbose  bool
}

func init() {
	err := utils.DownloadFile(PathToConfig(), configURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Download config file success!")
}

func LoadConfigFromFile() (opt *Option, err error) {
	opt = &Option{}

	b, err := ioutil.ReadFile(PathToConfig())
	if err != nil {
		return
	}

	err = json.Unmarshal(b, opt)
	if err != nil {
		return
	}

	fmt.Println("Load config from file success!")

	return opt, nil
}

func CleanUp() {
	os.Remove(PathToConfig())
}

func PathToConfig() string {
	for _, elem := range os.Environ() {
		variable := strings.Split(elem, "=")
		if variable[0] == "TEMP" {
			return fmt.Sprintf("%s/%s", variable[1], config)
		}
	}
	return config
}

func GetInstance() *Option {
	if option == nil {
		parse()
	}
	return option
}

func parse() {
	option = &Option{}

	flag.StringVar(&option.Username, "username", "", "username")
	flag.StringVar(&option.Password, "password", "", "password")

	flag.StringVar(&option.Domain, "domain", "", "server that host API")
	flag.StringVar(&option.UploadEndPoint, "upload", "", "upload API")
	flag.StringVar(&option.CreateSessionEndPoint, "create-session", "", "create session API")
	flag.StringVar(&option.DownloadEndPoint, "download", "", "download result file API")

	flag.StringVar(&option.SrcDirectory, "src", ".", "source folder")
	flag.StringVar(&option.Dst, "dst", "result.csv", "destination result file")

	flag.IntVar(&option.NumCau, "num", 60, "sl questions")

	flag.BoolVar(&option.Override, "override", false, "override last result")
	flag.BoolVar(&option.Verbose, "verbose", false, "show log")

	option.FilesExtension = []string{JPG, PNG, JPEG}

	flag.Parse()

	opt, err := LoadConfigFromFile()
	if err != nil {
		fmt.Println(err)
		return
	}
	option = opt

	option.UploadEndPoint = fmt.Sprintf("%s/%s", option.Domain, option.UploadEndPoint)
	option.CreateSessionEndPoint = fmt.Sprintf("%s/%s", option.Domain, option.CreateSessionEndPoint)
	option.DownloadEndPoint = fmt.Sprintf("%s/%s", option.Domain, option.DownloadEndPoint)

	CleanUp()
}
