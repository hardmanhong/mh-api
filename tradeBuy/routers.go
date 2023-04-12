package tradeBuy

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(router *gin.RouterGroup, db *gorm.DB) {
	// 将 Controller 结构体的实例作为指针传递给 Routers 函数
	dao := NewTradeDAO(db)
	controller := &TradeController{dao}
	// 使用指针接收者来调用 GetList 方法
	router.GET("/trade/buy", controller.GetList)
	router.GET("/trade/buy/:id", controller.GetItem)
	router.POST("/trade/buy", controller.Create)
	router.PUT("/trade/buy/:id", controller.Update)
	router.DELETE("/trade/buy/:id", controller.Delete)
}
