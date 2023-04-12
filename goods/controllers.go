package goods

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/utils"
)

type GoodsController struct {
	dao GoodsDAO
}

func (controller *GoodsController) GetList(ctx *gin.Context) {
	name := ctx.Query("name")
	page, pageSize := utils.GetPaginationParams(ctx)
	result, err := controller.dao.GetList(&ListQuery{
		Name:     name,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.ApiSuccessResponse(gin.H{"list": &result.List, "total": &result.Total})
	ctx.JSON(http.StatusOK, res)
}
func (controller *GoodsController) GetItem(ctx *gin.Context) {
	// 从 URL 参数中获取商品 ID
	var res gin.H
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res = utils.ApiErrorResponse(-1, "Invalid ID")
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// 调用 DAO 层获取商品
	goods, err := controller.dao.GetItem(id)
	if err != nil {
		res = utils.ApiErrorResponse(-1, "Failed to get item")
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// 返回结果
	if goods == nil {
		res = utils.ApiErrorResponse(-1, "Item not found")
		ctx.JSON(http.StatusNotFound, res)
	} else {
		res = utils.ApiSuccessResponse(&goods)
		ctx.JSON(http.StatusOK, res)
	}
}
func (controller *GoodsController) Create(ctx *gin.Context) {
	var goods Goods
	var res gin.H
	err := ctx.ShouldBindJSON(&goods)
	if err != nil {
		res = utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err = controller.dao.Create(&goods)
	if err != nil {
		res = utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	ctx.JSON(http.StatusOK, utils.ApiSuccessResponse(&goods))
}

func (controller *GoodsController) Update(ctx *gin.Context) {
	var req GoodsUpdate
	var res gin.H
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		res = utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res = utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// 检查记录是否存在
	exists, err := controller.dao.Exists(id)
	if err != nil {
		res = utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	if !exists {
		res = utils.ApiErrorResponse(404, "记录不存在")
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	if err := controller.dao.Update(id, &req); err != nil {
		res = utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	ctx.JSON(http.StatusOK, utils.ApiSuccessResponse(nil))
}
func (controller *GoodsController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ApiErrorResponse(-1, "Invalid id parameter"))
		return
	}

	err = controller.dao.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ApiErrorResponse(-1, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.ApiSuccessResponse(nil))
}
