package controller

import (
	"btools/pkg/service/export_file"
	"btools/pkg/service/parse_tools"
	"btools/pkg/service/parse_upload"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetTsList
func GetTsList(c *gin.Context) {
	formFile, _ := c.FormFile("file")
	uploadFile, err := formFile.Open()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	listData, err := parse_upload.ParseList(uploadFile)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	exportName := strconv.FormatInt(time.Now().Unix(), 10) + "_ts.txt"
	filePath, _ := export_file.CreateExportFile(exportName)
	for _, item := range listData {
		req, _ := parse_tools.AddReqItem(item)
		data, err := req.ParseM3u8()
		if err == nil {
			export_file.AppendExportFile(exportName, req.Host, data)
		}
	}
	//fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "ts_"+formFile.Filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(filePath)
}
