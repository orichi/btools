package parse_and_export

import (
	"btools/pkg/service/export_file"
	"btools/pkg/service/parse_tools"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	bufferSize = 2 << 20 // 缓存大小为1m
)

// Process
// 处理m3u8列表，然后写入文件，并返回文件的路径
func Process(urlList []string) (string, error) {
	exportName := strconv.FormatInt(time.Now().Unix(), 10) + "_ts.txt"
	filePath, err := export_file.CreateExportFile(exportName)
	if err != nil {
		return "", err
	}

	var tsListChan = make(chan []string, 5000)
	var wg = new(sync.WaitGroup)
	// 把url发送到请求chan里
	wg.Add(len(urlList))
	go func(group *sync.WaitGroup) {
		for _, item := range urlList {
			wg.Done()
			reqItem, err := parse_tools.AddReqItem(item)
			if err == nil {
				if tsList, err := reqItem.ParseM3u8(); err != nil {
					fmt.Println("err", err)
				} else {
					tsListChan <- tsList
				}
			}
		}
		close(tsListChan)
	}(wg)
	wg.Wait()

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0666)
	defer f.Close()
	buf := bufio.NewWriterSize(f, bufferSize)

	for {
		tsArray, ok := <-tsListChan
		if !ok {
			break
		}
		for _, item := range tsArray {
			_, err = buf.WriteString(item)
			if err != nil {
				continue
			}
		}
	}
	buf.Flush()

	return filePath, nil
}
