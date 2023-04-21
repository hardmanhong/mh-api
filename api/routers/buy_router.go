package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/controllers"
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/services"
	"gorm.io/gorm"
)

func NewBuyRouter(router *gin.RouterGroup, db *gorm.DB) {
	buyDAO := dao.NewBuyDAO(db)
	buyService := services.NewBuyService(buyDAO)
	buyController := controllers.NewBuyController(buyService)
	router.GET("/buy", buyController.GetList)
	router.GET("/buy/:id", buyController.GetItem)
	router.POST("/buy", buyController.Create)
	router.PUT("/buy/:id", buyController.Update)
	router.DELETE("/buy/:id", buyController.Delete)
}
