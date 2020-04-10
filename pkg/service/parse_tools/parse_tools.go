package parse_tools

import (
	"btools/pkg/configure"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/viki-org/dnscache"
)

var (
	ErrReqFail  = errors.New("请求失败")
	ErrNotFound = errors.New("未找到")
	ErrParseUrl = errors.New("URL解析失败")
	ctx, _      = context.WithTimeout(nil, 5*time.Second)
	dnsCache    = dnscache.New(time.Minute * 5)
	transport   = &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost: 100,
		DialContext: func(ctx context.Context, network string, address string) (net.Conn, error) {
			separator := strings.LastIndex(address, ":")
			ip, _ := dnsCache.FetchOneString(address[:separator])
			var a = new(net.Dialer)
			a.Timeout = 3 * time.Second    // 连接超时
			a.KeepAlive = 60 * time.Second //
			return a.DialContext(ctx, "tcp", ip+address[separator:])
		},
		IdleConnTimeout: 90,
	}
	// 超时时间自行设置
	httpClient = http.Client{
		Timeout:   time.Duration(configure.Conf.HttpWaitPeriod) * time.Second,
		Transport: transport,
	}

	regexM3u8 = regexp.MustCompile("\\.m3u8")
	regexTs   = regexp.MustCompile("\\.ts")
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

func (reqItem *ReqItem) ParseM3u8() ([]string, error) {

loop:
	req := &http.Request{
		URL: reqItem.ReqURL,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return nil, ErrReqFail
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("请求失败:", resp.StatusCode)
		return nil, ErrNotFound
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	tsList, err, retryURL := ParseTsItems(reqItem.Host, resp.StatusCode, respBytes)
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
func ParseTsItems(host string, state int, Data []byte) ([]string, error, string) {
	if state != 200 {
		return nil, ErrReqFail, ""
	}
	var tsItems []string
	buffer := bytes.NewBuffer(Data)
	for {
		line, err := buffer.ReadString('\n')

		if err != nil {
			break
		}
		lenLine := len(line)
		if lenLine > 4 && line[0:4] == "#EXT" {
			continue
		}

		if regexTs.MatchString(line) {
			tsItems = append(tsItems, host+line)
			continue
		}

		if regexM3u8.MatchString(line) {
			return nil, nil, strings.TrimSpace(line)
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
