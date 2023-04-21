package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/services"
)

type statisticsController struct {
	service services.StatisticsService
}

func NewStatisticsController(service services.StatisticsService) *statisticsController {
	return &statisticsController{service: service}
}

func (controller *statisticsController) GetStatistics(ctx *gin.Context) {
	dType := ctx.Query("type")
	res := controller.service.GetStatistics(dType)
	ctx.JSON(http.StatusOK, res)
}

func (controller *statisticsController) GetTotalProfit(ctx *gin.Context) {
	res := controller.service.GetTotalProfit()
	ctx.JSON(http.StatusOK, res)
}
