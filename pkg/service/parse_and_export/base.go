package parse_and_export

import (
	"btools/pkg/service/export_file"
	"btools/pkg/service/parse_tools"
	"bufio"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
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

	var listChan = make(chan string, 20)
	var tsChan = make(chan string, 5000)
	// 把url发送到请求chan里
	go func() {
		for i := 0; i < len(urlList); i++ {
			listChan <- urlList[i]
		}
		close(listChan)
	}()

	go func() {
		var wg = new(sync.WaitGroup)
		for url := range listChan {
			wg.Add(1)
			go func(url string) {
				reqItem, err := parse_tools.AddReqItem(url)

				if err == nil {
					if tsList, err := reqItem.ParseM3u8(); err != nil {
						logrus.Error(err.Error())
					} else {
						for _, item := range tsList {
							tsChan <- item
						}
					}
					wg.Done()
				}

			}(url)
		}
		wg.Wait()
		close(tsChan)
	}()
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0666)
	defer f.Close()
	buf := bufio.NewWriterSize(f, bufferSize)

	for {
		tsItem, ok := <-tsChan
		if !ok {
			break
		}
		buf.WriteString(tsItem)
	}
	buf.Flush()

	return filePath, nil
}
