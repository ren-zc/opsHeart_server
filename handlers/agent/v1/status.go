package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opsHeart_server/common"
	"opsHeart_server/db"
	"opsHeart_server/logger"
	"opsHeart_server/service/agent"
)

func HandleStatus(c *gin.Context) {
	UUID, ok := c.Get(common.UUIDKey)
	if !ok {
		logger.HbsLog.Errorf("path=status;action=get uuid;ip=%s", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "can not get uuid",
		})
		return
	}
	agentUuid := UUID.(string)

	var q agent.Agent
	err := db.DB.Find(&q, "uuid = ?", agentUuid).Error
	if err != nil {
		logger.RigisterLog.Error("uuid=%s, remote_ip=%s, action=query db, err=%s", agentUuid,
			c.Request.RemoteAddr, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": q.Status,
	})
}
