package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/services"
	"github.com/hardmanhong/api/utils"
)

type CharacterController struct {
	service services.CharacterService
}

func NewCharacterController(service services.CharacterService) *CharacterController {
	return &CharacterController{service: service}
}

func (controller *CharacterController) GetList(ctx *gin.Context) {
	res := controller.service.GetList()
	ctx.JSON(http.StatusOK, res)
}

func (controller *CharacterController) GetItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		res := utils.ApiErrorResponse(-1, "Invalid ID")
		ctx.JSON(http.StatusOK, res)
		return
	}

	// 调用 DAO 层获取商品
	res := controller.service.GetItem(uint32(id))
	// 返回结果
	ctx.JSON(http.StatusOK, res)
}

func (controller *CharacterController) Create(ctx *gin.Context) {
	var item models.Character
	err := ctx.ShouldBindJSON(&item)
	jsonData, _ := json.Marshal(item)
	println("query", string(jsonData))

	if err != nil {
		res := utils.ApiErrorResponse(-1, err.Error())
		ctx.JSON(http.StatusOK, res)
		return
	}

	res := controller.service.Create(&item)
	ctx.JSON(http.StatusOK, res)
}

func (controller *CharacterController) Update(ctx *gin.Context) {
	var req models.Character
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

func (controller *CharacterController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ApiErrorResponse(-1, "Invalid id parameter"))
		return
	}
	res := controller.service.Delete(uint32(id))
	ctx.JSON(http.StatusOK, res)
}
