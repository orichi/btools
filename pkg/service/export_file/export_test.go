package export_file

import (
	"fmt"
	"testing"
)

func TestExportFile(t *testing.T) {
	data := [][]byte{[]byte("/asdasd/11.ts\n")}
	exportFilePath, err := ExportFile("aaa.txt", "http://www.baidu.com", data)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(exportFilePath)
}
