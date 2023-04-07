package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hardmanhong/goods"
	"github.com/hardmanhong/utils"
)

func main() {
	// 1.创建路由
	r := gin.Default()
	api := r.Group("/api")
	goods.Routers(api)

	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/", func(c *gin.Context) {
		res := utils.FormatResult(200, gin.H{"foo": "bar"}, "ok")
		c.JSON(http.StatusOK, res)
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run()
}
