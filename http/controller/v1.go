package controller

import (
	"btools/pkg/configure"
	"btools/pkg/service/parse_and_export"
	"btools/pkg/service/parse_upload"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// GetTsList
func GetTsList(c *gin.Context) {
	formFile, err := c.FormFile("file")
	if err != nil {
		logrus.Error(err)
		c.JSON(500, "未发现上传文件")
		return
	}

	uploadFile, err := formFile.Open()
	if err != nil {
		logrus.Error(err)
		c.JSON(500, "上传文件打开失败")
		return
	}

	// 解析上传文件
	listData, err := parse_upload.ParseList(uploadFile)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, "列表解析失败")
		return
	}

	// 请求并导出到文件名

	var nums = len(listData)
	var toBeDealSize = nums
	if toBeDealSize == 0 {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":    "请上传文件",
			"path":     "/ts_list",
			"error":    "待处理记录为空，请重新选择文件",
			"max_line": configure.Conf.MaxLine,
		})
		return
	}
	if configure.Conf.MaxLine < nums {
		toBeDealSize = configure.Conf.MaxLine
	}

	data := parse_and_export.Process(listData[:toBeDealSize])

	//fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "ts_"+formFile.Filename))
	c.Data(http.StatusOK, "application/octet-stream", data)
}
