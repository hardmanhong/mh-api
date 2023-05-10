package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/services"
	"github.com/hardmanhong/api/utils"
)

type AccountController struct {
	service services.AccountService
}

func NewAccountController(service services.AccountService) *AccountController {
	return &AccountController{service: service}
}

func (controller *AccountController) GetList(ctx *gin.Context) {
	userName := ctx.Query("userName")
	userID := ctx.GetString("userID")
	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		userId = 0
	}
	res := controller.service.GetList(userId, &models.AccountListQuery{
		Name: userName,
	})
	ctx.JSON(http.StatusOK, res)
}

func (controller *AccountController) GetItem(ctx *gin.Context) {
	// 从 URL 参数中获取商品 ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		res := utils.ApiErrorResponse(-1, "ID 不存在")
		ctx.JSON(http.StatusOK, res)
		return
	}
	res := controller.service.GetItem(uint32(id))
	ctx.JSON(http.StatusOK, res)
}

func (controller *AccountController) Create(ctx *gin.Context) {
	var account models.Account
	err := ctx.ShouldBindJSON(&account)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}
	userID := ctx.GetString("userID")
	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}

	res := controller.service.Create(userId, &account)
	ctx.JSON(http.StatusOK, res)
}

func (controller *AccountController) Update(ctx *gin.Context) {
	var req models.AccountUpdate
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

	res := controller.service.Update(uint32(id), &req)
	ctx.JSON(http.StatusOK, res)
}

func (controller *AccountController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ApiErrorResponse(-1, "Invalid id parameter"))
		return
	}

	res := controller.service.Delete(uint32(id))
	ctx.JSON(http.StatusOK, res)
}
