package controller

import (
	"btools/pkg/service/parse_and_export"
	"btools/pkg/service/parse_upload"
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetTsList
func GetTsList(c *gin.Context) {
	formFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	uploadFile, err := formFile.Open()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	// 解析上传文件
	listData, err := parse_upload.ParseList(uploadFile)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	// 请求并导出到文件名
	filePath, _ := parse_and_export.Process(listData)

	//fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "ts_"+formFile.Filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(filePath)
}
