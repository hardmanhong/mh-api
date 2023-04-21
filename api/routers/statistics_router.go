package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/controllers"
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/services"
	"gorm.io/gorm"
)

func NewStatisticsRouter(router *gin.RouterGroup, db *gorm.DB) {
	dao := dao.NewStatisticsDAO(db)
	service := services.NewStatisticsService(dao)
	controller := controllers.NewStatisticsController(service)
	router.GET("/statistics", controller.GetStatistics)
	router.GET("/statistics/totalProfit", controller.GetTotalProfit)
}
