package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/controllers"
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/services"
	"gorm.io/gorm"
)

func NewSellRouter(router *gin.RouterGroup, db *gorm.DB) {
	buyDAO := dao.NewBuyDAO(db)
	sellDAO := dao.NewSellDAO(db)
	sellService := services.NewSellService(sellDAO, buyDAO)
	sellController := controllers.NewSellController(sellService)
	router.GET("/sell/:id", sellController.GetItem)
	router.POST("/sell", sellController.Create)
	router.PUT("/sell/:id", sellController.Update)
	router.DELETE("/sell/:id", sellController.Delete)
}
