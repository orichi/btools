package cmd

import (
	"btools/pkg/configure"
	"btools/pkg/service/parse_and_export"
	"btools/pkg/service/parse_upload"
	"errors"
	"fmt"
	"os"
)

var (
	ErrSourceFile = errors.New("源文件解析失败")
)

// Start 解析文件的方式去解析ts文件列表
func Start(sourceFile string) error {
	listFile, err := os.Open(sourceFile)
	if err != nil {
		return err
	}

	listData, err := parse_upload.ParseList(listFile)
	if err != nil {
		return err
	}
	var nums = len(listData)
	var toBeDealSize = nums
	if configure.Conf.MaxLine < nums {
		toBeDealSize = configure.Conf.MaxLine
	}

	filePath, _ := parse_and_export.ProcessFile(listData[:toBeDealSize])
	fmt.Println("写入文件:", filePath)
	return nil
}
