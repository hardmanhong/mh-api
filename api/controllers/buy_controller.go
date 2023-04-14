package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/services"
	"github.com/hardmanhong/api/utils"
)

type BuyController struct {
	service services.BuyService
}

func NewBuyController(service services.BuyService) *BuyController {
	return &BuyController{service: service}
}

func (controller *BuyController) GetList(ctx *gin.Context) {
	page, pageSize := utils.GetPaginationParams(ctx)
	// 解析日期参数
	createdAtFrom, createdAtTo := utils.ParseDate([2]string{"createdAtFrom", "createdAtTo"}, ctx.Request.URL.Query())
	// 获取数组参数 buyIds
	buyIdsQuery, exists := ctx.Request.URL.Query()["buyIds[]"]
	var buyIds []uint64
	if exists {
		for _, idStr := range buyIdsQuery {
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				buyIds = nil
				break
			}
			buyIds = append(buyIds, id)
		}
	}
	result, err := controller.service.GetList(&models.BuyListQuery{
		CreatedAtFrom: createdAtFrom,
		CreatedAtTo:   createdAtTo,
		GoodsIDs:      buyIds,
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
	res := utils.ApiSuccessResponse(&result)
	ctx.JSON(http.StatusOK, res)
}

func (controller *BuyController) GetItem(ctx *gin.Context) {
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

func (controller *BuyController) Create(ctx *gin.Context) {
	var buy models.Buy
	err := ctx.ShouldBindJSON(&buy)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}

	b, err := controller.service.Create(&buy)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}
	ctx.JSON(http.StatusOK, utils.ApiSuccessResponse(&b))
}

func (controller *BuyController) Update(ctx *gin.Context) {
	var req models.BuyUpdate
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

func (controller *BuyController) Delete(ctx *gin.Context) {
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
