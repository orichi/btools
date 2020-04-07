package parse_tools

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	ErrReqFail = errors.New("请求失败")

	ErrNotFound = errors.New("未找到")
	ErrParseUrl = errors.New("URL解析失败")

	httpClient = http.Client{Timeout: 10 * time.Second, Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	regexM3u8 = regexp.MustCompile(".m3u8")
	regexTs   = regexp.MustCompile(".ts")
)

type ReqItem struct {
	Host    string   // 带有https?
	Domain  string   // 域名
	ReqURL  *url.URL // 请求路径全路径
	Retry   bool     // 是否是第二次请求
	Success bool     // 请求是否成功
}

func AddReqItem(reqURL string) (*ReqItem, error) {
	i, err := FindDelimiterIndex(reqURL)
	if err != nil {
		return nil, err
	}

	// 长度应该 >= len(host)+len('/1.ts')
	if len(reqURL)-i <= 2 {
		return nil, ErrParseUrl
	}

	URL, _ := url.Parse(reqURL)
	item := &ReqItem{
		Host:   reqURL[0:i],
		ReqURL: URL,
	}
	return item, nil
}

func (reqItem *ReqItem) ParseM3u8() ([][]byte, error) {

loop:
	req := &http.Request{
		URL: reqItem.ReqURL,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	tsList, err, retryURL := ParseTsItems(resp.StatusCode, respBytes)
	if retryURL != "" {
		URL, err := url.Parse(reqItem.Host + retryURL)
		if err != nil {
			return nil, err
		}
		reqItem.ReqURL = URL
		reqItem.Retry = true
		goto loop
	}
	if err != nil {
		return nil, err
	}

	return tsList, nil
}

// ParseTsItems 解析ts列表
func ParseTsItems(state int, Data []byte) ([][]byte, error, string) {
	if state != 200 {
		return nil, ErrReqFail, ""
	}
	var tsItems = [][]byte{}
	buffer := bytes.NewBuffer(Data)
	for {
		line, err := buffer.ReadBytes('\n')

		if err != nil {
			break
		}
		lenLine := len(line)
		if lenLine > 4 && string(line[0:4]) == "#EXT" {
			continue
		}

		if regexTs.Match(line) {
			tsItems = append(tsItems, line[0:lenLine])
			continue
		}

		if regexM3u8.Match(line) {
			return nil, nil, strings.TrimSpace(string(line))
		}
	}

	return tsItems, nil, ""
}

func FindDelimiterIndex(url string) (int, error) {
	for i := 0; i < len(url); i++ {
		if i > 1 && url[i] == '/' && url[i+1] != '/' && url[i-1] != '/' {
			return i, nil
			break
		}

		for url[i+1] != '/' {
			i++
			break
		}
	}

	return 0, ErrNotFound
}
