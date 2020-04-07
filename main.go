package main

import (
	"btools/cmd"
	"btools/http"
	"fmt"
	"os"
	"regexp"
)

var (
	action   string
	addr     string
	showHelp bool
	version  = "0.1"
)

var (
	usage = `Usage:
  服务可以通过http方式、文件操作两种方式使用
    1. 如果需要启动http，通过http来操作的话
      命令：main.exe http :10001
      
           可执行文件 服务方式 监听端口号
      如果不需要外网访问，请在监听端口号前加 127.0.0.1, eg: 127.0.0.1:10001
      如果是windows，请用localhost代替127.0.0.1

      浏览器访问 127.0.0.1:10001/
      上传文本文件（utf8编码格式）

    2. 	直接操作文件的方式
      命令：main.exe file filename.txt
      
           可执行文件 服务方式 待处理的文件名
      会输出文件到 可执行文件目录里的public目录下
  
  Note: 不过是文件方式，还是http服务方式，txt文件都需要是utf8编码 
`
)

func main() {
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if len(os.Args) <= 1 {
		fmt.Println("请加上参数调用\n \t[-v|-h|http|file]")
		os.Exit(0)
	}

	runMethod := os.Args[1]
	var extraOptions string
	if len(os.Args) >= 3 {
		extraOptions = os.Args[2]
	}

	if runMethod == "http" {
		if match, err := regexp.MatchString("^([0-9a-z].)*?:[0-9]+$", extraOptions); err != nil || !match {
			fmt.Println("参数有误")
		} else {
			http.Run(extraOptions)
		}

	}

	if runMethod == "file" {
		if extraOptions == "" {
			fmt.Println("参数有误")
		} else {

			cmd.Start(extraOptions)
		}
	}

	if runMethod == "-v" {
		fmt.Println(version)
		os.Exit(0)
	}

	if runMethod == "-h" {
		fmt.Print(usage)
	}

	os.Exit(0)
}
