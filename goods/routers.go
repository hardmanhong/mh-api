package goods

import "github.com/gin-gonic/gin"

func Routers(router *gin.RouterGroup) {
	// 将 Controller 结构体的实例作为指针传递给 Routers 函数
	controller := &Controller{}
	// 使用指针接收者来调用 GetList 方法
	router.POST("/user", controller.GetList)
}
