package goods

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(router *gin.RouterGroup, db *gorm.DB) {
	// 将 Controller 结构体的实例作为指针传递给 Routers 函数
	dao := NewGoodsDAO(db)
	controller := &GoodsController{dao}
	// 使用指针接收者来调用 GetList 方法
	router.GET("/goods", controller.GetList)
	router.GET("/goods/:id", controller.GetItem)
	router.POST("/goods", controller.Create)
	router.PUT("/goods/:id", controller.Update)
	router.DELETE("/goods/:id", controller.Delete)
}
