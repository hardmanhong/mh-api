package tradeSell

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/utils"
)

type TradeController struct {
	dao TradeSellDAO
}

func (controller *TradeController) GetItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	result, err := controller.dao.GetItem(id)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	if result == nil {
		res := utils.ApiErrorResponse(-1, "record not found")
		ctx.JSON(http.StatusNotFound, res)
		return
	}
	res := utils.ApiSuccessResponse(result)
	ctx.JSON(http.StatusOK, res)
}

func (controller *TradeController) Create(ctx *gin.Context) {
	var trade TradeSell
	err := ctx.ShouldBindJSON(&trade)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	t, err := controller.dao.Create(&trade)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.ApiSuccessResponse(t)
	ctx.JSON(http.StatusCreated, res)
}

func (controller *TradeController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	var tradeUpdate TradeSellUpdate
	if err := ctx.ShouldBindJSON(&tradeUpdate); err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	if err := controller.dao.Update(id, &tradeUpdate); err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.ApiSuccessResponse(nil)
	ctx.JSON(http.StatusOK, res)
}

func (controller *TradeController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	if err := controller.dao.Delete(id); err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.ApiSuccessResponse(nil)
	ctx.JSON(http.StatusOK, res)
}
