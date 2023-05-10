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
	userID := ctx.GetString("userID")
	page, pageSize := utils.GetPaginationParams(ctx)
	// 解析日期参数
	createdAtFrom, createdAtTo := utils.ParseDate([2]string{"createdAtFrom", "createdAtTo"}, ctx.Request.URL.Query())
	// 获取数组参数 goodsIds
	goodsIdsQuery, exists := ctx.Request.URL.Query()["goodsIds[]"]
	var goodsIds []uint64
	if exists {
		for _, idStr := range goodsIdsQuery {
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				goodsIds = nil
				break
			}
			goodsIds = append(goodsIds, id)
		}
	}
	inventorySorter := ctx.Query("inventorySorter")
	hasSoldSorter := ctx.Query("hasSoldSorter")
	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		userId = 0
	}
	res := controller.service.GetList(userId, &models.BuyListQuery{
		CreatedAtFrom:   createdAtFrom,
		CreatedAtTo:     createdAtTo,
		GoodsIDs:        goodsIds,
		InventorySorter: inventorySorter,
		HasSoldSorter:   hasSoldSorter,
		PaginationQuery: models.PaginationQuery{
			Page:     page,
			PageSize: pageSize,
		},
	})
	ctx.JSON(http.StatusOK, res)
}

func (controller *BuyController) GetItem(ctx *gin.Context) {
	// 从 URL 参数中获取商品 ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		res := utils.ApiErrorResponse(-1, "Invalid ID")
		ctx.JSON(http.StatusOK, res)
		return
	}

	res := controller.service.GetItem(id)
	ctx.JSON(http.StatusOK, res)
}

func (controller *BuyController) Create(ctx *gin.Context) {
	var buy models.Buy
	err := ctx.ShouldBindJSON(&buy)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}
	userID := ctx.GetString("userID")
	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		userId = 0
	}
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}

	res := controller.service.Create(userId, &buy)
	ctx.JSON(http.StatusOK, res)
}

func (controller *BuyController) Update(ctx *gin.Context) {
	var req models.BuyUpdate
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

func (controller *BuyController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ApiErrorResponse(-1, "参数无效"))
		return
	}
	res := controller.service.Delete(id)
	ctx.JSON(http.StatusOK, res)
}
