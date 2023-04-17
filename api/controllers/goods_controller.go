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
	res := controller.service.GetList(&models.GoodsListQuery{
		Name: name,
		PaginationQuery: models.PaginationQuery{
			Page:     page,
			PageSize: pageSize,
		},
	})
	ctx.JSON(http.StatusOK, res)
}

func (controller *GoodsController) GetItem(ctx *gin.Context) {
	// 从 URL 参数中获取商品 ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		res := utils.ApiErrorResponse(-1, "ID 不存在")
		ctx.JSON(http.StatusOK, res)
		return
	}
	res := controller.service.GetItem(id)
	ctx.JSON(http.StatusOK, res)
}

func (controller *GoodsController) Create(ctx *gin.Context) {
	var goods models.Goods
	err := ctx.ShouldBindJSON(&goods)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}

	res := controller.service.Create(&goods)
	ctx.JSON(http.StatusOK, res)
}

func (controller *GoodsController) Update(ctx *gin.Context) {
	var req models.GoodsUpdate
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

func (controller *GoodsController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ApiErrorResponse(-1, "Invalid id parameter"))
		return
	}

	res := controller.service.Delete(id)
	ctx.JSON(http.StatusOK, res)
}
