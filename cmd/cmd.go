package cmd

import (
	"btools/pkg/service/export_file"
	"btools/pkg/service/parse_tools"
	"btools/pkg/service/parse_upload"
	"errors"
	"os"
	"strconv"
	"time"
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

	exportName := strconv.FormatInt(time.Now().Unix(), 10) + "_ts.txt"
	export_file.CreateExportFile(exportName)
	for _, item := range listData {
		req, _ := parse_tools.AddReqItem(item)
		data, err := req.ParseM3u8()
		if err == nil {
			export_file.AppendExportFile(exportName, req.Host, data)
		}
	}

	return nil
}
