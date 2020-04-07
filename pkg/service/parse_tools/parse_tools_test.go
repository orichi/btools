package parse_tools

import (
	"fmt"
	"testing"
)

func TestfindDelimiterIndex(t *testing.T) {
	var url = "https://www.mzy2000.com:65/20200111/ynnzCsJB/index.m3u8"
	index, err := FindDelimiterIndex(url)
	if err != nil {
		t.Error(err, index)
	}
}

func TestAddReqItem(t *testing.T) {
	var url = "https://www.mzy2000.com:65/20200111/ynnzCsJB/index.m3u8"
	fmt.Println(url[0:26])
	item, err := AddReqItem(url)
	if err != nil {
		t.Error(err)
		return
	}
	if item.Host != "https://www.mzy2000.com:65" {
		t.Error("校验失败: ", item)
	}
}

func TestReqItem_ParseM3u8(t *testing.T) {
	reqItem, err := AddReqItem("https://www.mzy2000.com:65/20200111/ynnzCsJB/index.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	tsList, err := reqItem.ParseM3u8()
	if err != nil {
		t.Fatal(err)
	}
	if len(tsList) == 0 {
		t.Fatal("未获取到数据")
	}
}
