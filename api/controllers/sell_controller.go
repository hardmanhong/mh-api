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
	var res gin.H
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		res = utils.ApiErrorResponse(-1, "Invalid ID")
		ctx.JSON(http.StatusOK, res)
		return
	}

	// 调用 DAO 层获取商品
	buy, err := controller.service.GetItem(id)
	if err != nil {
		res = utils.ApiErrorResponse(-1, "Failed to get item")
		ctx.JSON(http.StatusOK, res)
		return
	}

	// 返回结果
	if buy == nil {
		res = utils.ApiErrorResponse(-1, "Item not found")
		ctx.JSON(http.StatusOK, res)
	} else {
		res = utils.ApiSuccessResponse(&buy)
		ctx.JSON(http.StatusOK, res)
	}
}

func (controller *SellController) Create(ctx *gin.Context) {
	var item models.Sell
	err := ctx.ShouldBindJSON(&item)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}

	b, err := controller.service.Create(&item)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}
	ctx.JSON(http.StatusOK, utils.ApiSuccessResponse(&b))
}

func (controller *SellController) Update(ctx *gin.Context) {
	var req models.SellUpdate
	var res gin.H
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		res = utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res = utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}

	// 检查记录是否存在
	exists, err := controller.service.Exists(id)
	if err != nil {
		res = utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}
	if !exists {
		res = utils.ApiErrorResponse(404, "记录不存在")
		ctx.JSON(http.StatusOK, res)
		return
	}

	if err := controller.service.Update(id, &req); err != nil {
		res = utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}
	ctx.JSON(http.StatusOK, utils.ApiSuccessResponse(nil))
}

func (controller *SellController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ApiErrorResponse(-1, "Invalid id parameter"))
		return
	}

	err = controller.service.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ApiErrorResponse(-1, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.ApiSuccessResponse(nil))
}
