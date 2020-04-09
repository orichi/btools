package http

import (
	"btools/http/controller"
	"btools/http/middleware"
	"btools/pkg/configure"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run(addr string) {
	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())
	r.MaxMultipartMemory = 10 << 20 // 最大附件大小10M
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":    "请上传文件",
			"path":     "/ts_list",
			"max_line": configure.Conf.MaxLine,
		})
	})

	r.POST("/ts_list", controller.GetTsList)
	err := r.Run(addr)
	if err != nil {
		log.Fatal("启动服务失败:", err.Error())
	}
}
