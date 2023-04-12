package tradeSell

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(router *gin.RouterGroup, db *gorm.DB) {
	// 将 Controller 结构体的实例作为指针传递给 Routers 函数
	dao := NewTradeDAO(db)
	controller := &TradeController{dao}
	// 使用指针接收者来调用 GetList 方法
	router.GET("/trade/sell/:id", controller.GetItem)
	router.POST("/trade/sell", controller.Create)
	router.PUT("/trade/sell/:id", controller.Update)
	router.DELETE("/trade/sell/:id", controller.Delete)
}
