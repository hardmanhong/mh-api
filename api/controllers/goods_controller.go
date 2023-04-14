package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/services"
	"github.com/hardmanhong/api/utils"
)

type GoodsController struct {
	service services.GoodsService
}

func NewGoodsController(service services.GoodsService) *GoodsController {
	return &GoodsController{service: service}
}

func (controller *GoodsController) GetList(ctx *gin.Context) {
	name := ctx.Query("name")
	page, pageSize := utils.GetPaginationParams(ctx)
	result, err := controller.service.GetList(&models.GoodsListQuery{
		Name: name,
		PaginationQuery: models.PaginationQuery{
			Page:     page,
			PageSize: pageSize,
		},
	})
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}
	res := utils.ApiSuccessResponse(result)
	ctx.JSON(http.StatusOK, res)
}

func (controller *GoodsController) GetItem(ctx *gin.Context) {
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
	goods, err := controller.service.GetItem(id)
	if err != nil {
		res = utils.ApiErrorResponse(-1, "Failed to get item")
		ctx.JSON(http.StatusOK, res)
		return
	}

	// 返回结果
	if goods == nil {
		res = utils.ApiErrorResponse(-1, "Item not found")
		ctx.JSON(http.StatusOK, res)
	} else {
		res = utils.ApiSuccessResponse(&goods)
		ctx.JSON(http.StatusOK, res)
	}
}

func (controller *GoodsController) Create(ctx *gin.Context) {
	var goods models.Goods
	var res gin.H
	err := ctx.ShouldBindJSON(&goods)
	if err != nil {
		res = utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}

	err = controller.service.Create(&goods)
	if err != nil {
		res = utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}
	ctx.JSON(http.StatusOK, utils.ApiSuccessResponse(&goods))
}

func (controller *GoodsController) Update(ctx *gin.Context) {
	var req models.GoodsUpdate
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

func (controller *GoodsController) Delete(ctx *gin.Context) {
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
