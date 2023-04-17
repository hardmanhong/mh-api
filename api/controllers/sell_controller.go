package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/services"
	"github.com/hardmanhong/api/utils"
)

type SellController struct {
	service services.SellService
}

func NewSellController(service services.SellService) *SellController {
	return &SellController{service: service}
}
func (controller *SellController) GetItem(ctx *gin.Context) {
	// 从 URL 参数中获取商品 ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		res := utils.ApiErrorResponse(-1, "Invalid ID")
		ctx.JSON(http.StatusOK, res)
		return
	}

	// 调用 DAO 层获取商品
	res := controller.service.GetItem(id)
	// 返回结果
	ctx.JSON(http.StatusOK, res)
}

func (controller *SellController) Create(ctx *gin.Context) {
	var item models.Sell
	err := ctx.ShouldBindJSON(&item)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}

	res := controller.service.Create(&item)
	ctx.JSON(http.StatusOK, res)
}

func (controller *SellController) Update(ctx *gin.Context) {
	var req models.SellUpdate
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}
	res := controller.service.Update(id, &req)
	ctx.JSON(http.StatusOK, res)
}

func (controller *SellController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ApiErrorResponse(-1, "Invalid id parameter"))
		return
	}
	res := controller.service.Delete(id)
	ctx.JSON(http.StatusOK, res)
}
