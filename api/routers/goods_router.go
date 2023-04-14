package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/controllers"
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/services"
	"gorm.io/gorm"
)

func NewGoodsRouter(router *gin.RouterGroup, db *gorm.DB) {
	goodsDAO := dao.NewGoodsDAO(db)
	goodsService := services.NewGoodsService(goodsDAO)
	goodsController := controllers.NewGoodsController(goodsService)
	router.GET("/goods", goodsController.GetList)
	router.GET("/goods/:id", goodsController.GetItem)
	router.POST("/goods", goodsController.Create)
	router.PUT("/goods/:id", goodsController.Update)
	router.DELETE("/goods/:id", goodsController.Delete)
}
