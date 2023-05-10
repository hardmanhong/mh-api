package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/controllers"
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/services"
	"gorm.io/gorm"
)

func NewCharacterRouter(router *gin.RouterGroup, db *gorm.DB) {
	characterDAO := dao.NewCharacterDAO(db)
	equipmentDAO := dao.NewEquipmentDAO(db)
	petDAO := dao.NewPetDAO(db)
	characterService := services.NewCharacterService(characterDAO, equipmentDAO, petDAO)
	characterController := controllers.NewCharacterController(characterService)
	router.GET("/character", characterController.GetList)
	router.GET("/character/:id", characterController.GetItem)
	router.POST("/character", characterController.Create)
	router.PUT("/character/:id", characterController.Update)
	router.DELETE("/character/:id", characterController.Delete)
}
