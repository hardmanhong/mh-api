package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/controllers"
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/services"
	"gorm.io/gorm"
)

func NewUserRouter(router *gin.RouterGroup, db *gorm.DB) {
	buyDAO := dao.NewUserDAO(db)
	buyService := services.NewUserService(buyDAO)
	buyController := controllers.NewUserController(buyService)
	router.POST("/signup", buyController.SignUp)
	router.POST("/login", buyController.Login)
}
