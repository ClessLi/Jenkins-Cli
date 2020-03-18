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
	//confData, readErr := readFile("F:/GO_Project/src/Jenkins-Cli/cmd/jk_cli/conf.yaml")
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
	// 代码测试
	//views := "test_view1,test_view2"
	//viewList = strings.Split(views, ",")

	jobs, err := jkClinet.GetJobsByViews(viewList)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//jobsInfo := []string{"Job_Name\tBuild_Times\n"}
	jobsInfo := "Job_Name\tBuild_Times\n"
	for _, job := range jobs {
		job, err = jkClinet.GetJob(job.Name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//fmt.Println(job.Name, job.LastBuild.Number)
		//jobsInfo = append(jobsInfo, fmt.Sprintf("%s\t%d\n", job.Name, job.LastBuild.Number))
		jobsInfo = jobsInfo + fmt.Sprintf("%s\t%d\n", job.Name, job.LastBuild.Number)
	}
	fmt.Println(jobsInfo)
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
