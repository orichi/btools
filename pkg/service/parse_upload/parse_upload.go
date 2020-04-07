package parse_upload

// 解析上传的文件

import (
	"bufio"
	"errors"
	"io"
	"regexp"
)

var (
	ErrParseUploadFile = errors.New("上传文件解析失败, 文件看起来没有有效内容")
)

// ParseList 从url列表文件中解析出来 m3u8的列表slice
func ParseList(uploadFile io.ReadCloser) ([]string, error) {
	buffers := bufio.NewReader(uploadFile)

	if buffers.Size() < 5 {
		return nil, ErrParseUploadFile
	}

	var items []string
	for {
		line, _, err := buffers.ReadLine()
		if err != nil {
			break
		}

		if matched, err := regexp.Match(".m3u8", line); err == nil && matched {
			items = append(items, string(line))
		}
	}

	return items, nil
}
