package parse_and_export

import (
	"btools/pkg/service/parse_tools"
	"bytes"
	"sync"

	"github.com/sirupsen/logrus"
)

// Process
// 处理m3u8列表，然后写入文件，并返回文件的路径
func Process(urlList []string) []byte {
	var bufferData = new(bytes.Buffer)

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

	for {
		tsItem, ok := <-tsChan
		if !ok {
			break
		}
		bufferData.WriteString(tsItem)
	}

	return bufferData.Bytes()
}
