package db

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func Init() {
	// 连接 MySQL 数据库
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4")
	if err != nil {
		fmt.Println("Failed to connect to MySQL:", err)
		return
	}
	defer db.Close() // 延迟关闭数据库连接

	// Gin 框架示例代码
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})
	router.Run(":8080")
}
