package http

import (
	"btools/http/controller"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run(addr string) {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "请上传文件",
			"path":  "/ts_list",
		})
	})

	r.POST("/ts_list", controller.GetTsList)
	err := r.Run(addr)
	if err != nil {
		log.Fatal("启动服务失败:", err.Error())
	}
}
