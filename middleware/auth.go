package middleware

import (
	"IShare/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		//id, err := utils.ParseToken(token)
		id, err := strconv.ParseUint(token, 0, 64)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"msg": "token错误",
			})
			c.Abort()
			return
		}
		if user, notFound := service.GetUserByID(id); notFound {
			c.JSON(http.StatusBadGateway, gin.H{
				"msg": "用户不存在 or id不合法",
			})
			c.Abort()
		} else {
			c.Set("user", user)
		}
	}
}
