package middlewares

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/utils"
	"gorm.io/gorm"
)

func Auth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for signup and login requests
		if strings.Contains(c.Request.URL.Path, "signup") ||
			strings.Contains(c.Request.URL.Path, "login") {
			c.Next()
			return
		}
		// 从请求头中获取 token
		token := c.GetHeader("token")
		if token == "" {
			c.AbortWithStatusJSON(-1, utils.ApiErrorResponse(-1, "invalid token"))
			return
		}
		find := &models.Token{}
		err := db.Model(&models.Token{}).Where("token = ?", token).First(find).Error

		if err != nil || find.ID == 0 {
			c.AbortWithStatusJSON(-1, utils.ApiErrorResponse(-1, "invalid token"))
			return
		}
		if time.Now().After(find.ExpireAt) {
			c.AbortWithStatusJSON(-1, utils.ApiErrorResponse(-1, "凭证已过期，请重新登录"))
			return
		}

		// 将用户 id 信息放入 context 中
		fmt.Println("find.UserId", find.UserId)
		c.Set("userID", strconv.FormatUint(find.UserId, 10))

		// 继续处理请求
		c.Next()
	}
}
