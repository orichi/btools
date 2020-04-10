package parse_and_export

import (
	"btools/pkg/service/export_file"
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	bufferSize = 2 << 20 // 缓存大小为1m
)

// Process
// 处理m3u8列表，然后写入文件，并返回文件的路径
func ProcessFile(urlList []string) (string, error) {
	exportName := strconv.FormatInt(time.Now().Unix(), 10) + "_ts.txt"
	filePath, err := export_file.CreateExportFile(exportName)
	if err != nil {
		return "", err
	}

	byteData := Process(urlList)
	if len(byteData) == 0 {
		return "", errors.New("未获取到数据")
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	f.Write(byteData)

	return filePath, nil
}
