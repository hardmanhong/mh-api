package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/mh?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database!")
	// dbSQL, err := db.DB()
	// defer dbSQL.Close() // 在函数结束时自动关闭数据库连接
	return db
}
