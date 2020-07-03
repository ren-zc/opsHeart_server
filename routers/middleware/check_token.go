package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opsHeart_server/common"
	"opsHeart_server/logger"
	"opsHeart_server/service/agent"
)

func TokenChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		agentUuid, token, ok := c.Request.BasicAuth()
		if !ok {
			logger.ServerLog.Errorf("action=auth check;uuid=%s;path=%s;client=%s;err=basic auth false",
				agentUuid, c.FullPath(), c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "auth required",
			})
			c.Abort()
			return
		}

		var q agent.Agent
		// or query dat form redis
		err := q.QueryByUUID(agentUuid)
		if err != nil {
			logger.ServerLog.Errorf("action=auth check;uuid=%s;path=%s;client=%s;err=%s",
				agentUuid, c.FullPath(), c.ClientIP(), err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		if token != q.Token {
			logger.ServerLog.Errorf("action=auth check;uuid=%s;path=%s;client=%s;err=wrong token",
				agentUuid, c.FullPath(), c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "auth required",
			})
			c.Abort()
			return
		}

		c.Set(common.UUIDKey, agentUuid)
		c.Next()
	}
}
