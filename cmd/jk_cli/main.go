package main

import (
	"flag"
	"fmt"
	"github.com/ClessLi/golang-jenkins"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
	"strings"
)

var (
	confPath = flag.String("f", "", "jenkins `conf`.y(a)ml path.")
	views    = flag.String("v", "", "find jenkins job names with `views`, if you specify the view(s).")
	viewList []string
)

type jkConf struct {
	BaseUrl  string `yaml:"baseUrl"`
	User     string `yaml:"username"`
	ApiToken string `yaml:"token"`
}

func main() {
	flag.Parse()

	isExist, pathErr := PathExists(*confPath)
	if !isExist {
		if pathErr != nil {
			fmt.Println("The logfile", *confPath, "is not found.")
		} else {
			fmt.Println("Unkown error of the logfile.")
		}
		flag.Usage()
		os.Exit(1)
	}

	confData, readErr := readFile(*confPath)
	if readErr != nil {
		fmt.Println(readErr)
		os.Exit(1)
	}
	conf := &jkConf{}
	ymlErr := yaml.Unmarshal(confData, &conf)
	if ymlErr != nil {
		fmt.Println(ymlErr)
		os.Exit(1)
	}

	jkClinet := gojenkins.NewJenkins(&gojenkins.Auth{
		Username: conf.User,
		ApiToken: conf.ApiToken,
	}, conf.BaseUrl)

	if *views != "" {
		viewList = strings.Split(*views, ",")
	}

	jobs, err := jkClinet.GetJobsByViews(viewList)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, job := range jobs {
		fmt.Println(job.Name)
	}
}

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fd, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, err
	} else {
		return false, nil
	}
}
