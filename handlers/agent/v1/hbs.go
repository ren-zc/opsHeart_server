package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opsHeart/common"
	"opsHeart/logger"
	"opsHeart/service/agent"
)

func HandleHbs(c *gin.Context) {
	var agentIP string
	agentIP = c.ClientIP()

	UUID, ok := c.Get(common.UUIDKey)
	if !ok {
		logger.HbsLog.Errorf("path=hbs;action=get uuid;ip=%s", agentIP)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "can not get uuid",
		})
		return
	}

	agentUUID := UUID.(string)

	// save hbs time to db or redis, and log
	a := &agent.Agent{
		UUID: agentUUID,
	}
	if err := a.UpdateHbs(); err != nil {
		logger.HbsLog.Errorf("action=hbs;do=update hbs;uuid=%s;ip=%s;err=%s", agentUUID, agentIP, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	logger.HbsLog.Info("action=hbs;uuid=%s;ip=%s", agentUUID, agentIP)

	//response
	c.JSON(http.StatusOK, gin.H{
		"msg": "status",
	})
}
