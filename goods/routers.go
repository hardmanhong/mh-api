package goods

import "github.com/gin-gonic/gin"

func Routers(router *gin.RouterGroup) {
	// 将 Controller 结构体的实例作为指针传递给 Routers 函数
	dao := &GoodsDAOImpl{}
	controller := &GoodsController{goodsDao: dao}
	// 使用指针接收者来调用 GetList 方法
	router.GET("/goods", controller.GetList)
	router.GET("/goods/:id", controller.GetItem)
	router.POST("/goods", controller.Create)
	router.PUT("/goods/:id", controller.Update)
	router.DELETE("/goods/:id", controller.Delete)
}
