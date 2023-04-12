package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hardmanhong/database"
	"github.com/hardmanhong/goods"
	"github.com/hardmanhong/tradeBuy"
	"github.com/hardmanhong/tradeSell"
)

func main() {
	// 创建数据库连接
	db := database.NewDB()

	// 1.创建路由
	r := gin.Default()

	// 2.设置 CORS 头信息
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "token"}

	r.Use(cors.New(config))

	// 3.注册路由
	api := r.Group("/api")
	goods.Routers(api, db)
	tradeBuy.Routers(api, db)
	tradeSell.Routers(api, db)

	// 4.运行服务
	r.Run(":9000")
}
