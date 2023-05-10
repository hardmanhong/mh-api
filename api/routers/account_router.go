package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/controllers"
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/services"
	"gorm.io/gorm"
)

func NewAccountRouter(router *gin.RouterGroup, db *gorm.DB) {
	accountDAO := dao.NewAccountDAO(db)
	accountService := services.NewAccountService(accountDAO)
	accountController := controllers.NewAccountController(accountService)
	router.GET("/account", accountController.GetList)
	router.GET("/account/:id", accountController.GetItem)
	router.POST("/account", accountController.Create)
	router.PUT("/account/:id", accountController.Update)
	router.DELETE("/account/:id", accountController.Delete)
}
