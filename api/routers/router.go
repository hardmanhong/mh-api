package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(router *gin.RouterGroup, db *gorm.DB) {
	NewUserRouter(router, db)
	NewGoodsRouter(router, db)
	NewBuyRouter(router, db)
	NewSellRouter(router, db)
	NewStatisticsRouter(router, db)
	NewAccountRouter(router, db)
	NewCharacterRouter(router, db)
}
