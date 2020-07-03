package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opsHeart_server/common"
	"opsHeart_server/conf"
)

func FrontToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		frontToken := c.GetHeader(common.FrontTokenKey)
		if frontToken == "" || frontToken != conf.GetFrontToken() {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "token required.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
