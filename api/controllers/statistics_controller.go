package controllers

import (
	"net/http"
	"strconv"

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
	userID := ctx.GetString("userID")
	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		userId = 0
	}
	dType := ctx.Query("type")
	res := controller.service.GetStatistics(userId, dType)
	ctx.JSON(http.StatusOK, res)
}

func (controller *statisticsController) GetTotalProfit(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		userId = 0
	}
	res := controller.service.GetTotalProfit(userId)
	ctx.JSON(http.StatusOK, res)
}
