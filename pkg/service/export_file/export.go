package export_file

// 导出ts列表到文件中，host和对应的列表byte slice
import (
	"fmt"
	"log"
	"os"
)

const (
	Folder = "public/"
)

// ExportFile 把 host和data导出到filename
func ExportFile(filename string, host string, data [][]byte) (string, error) {
	_, err := os.Stat(Folder)
	if err != nil {
		log.Println(err)
		os.Mkdir(Folder, 0755)
		err = nil
	}

	f, err := os.Create(Folder + filename)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	for _, bytesItem := range data {
		fmt.Fprint(f, host)
		fmt.Fprint(f, string(bytesItem))
	}

	return Folder + filename, nil
}

// CreateExportFile 创建导出的文件名
func CreateExportFile(filename string) (string, error) {
	_, err := os.Stat(Folder)
	if err != nil {
		err = os.Mkdir(Folder, 0755)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.Create(Folder + filename)
	defer f.Close()

	if err != nil {
		return "", err
	}

	return Folder + filename, nil
}

// AppendExportFile 添加内容到导出文件
func AppendExportFile(filename string, data []string) error {
	f, err := os.OpenFile(Folder+filename, os.O_APPEND|os.O_WRONLY, 0666)
	defer f.Close()

	if err != nil {
		return err
	}

	for _, bytesItem := range data {
		fmt.Fprint(f, string(bytesItem))
	}

	return nil
}
