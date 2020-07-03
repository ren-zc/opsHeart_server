package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opsHeart_server/logger"
	"opsHeart_server/service/agent"
	"opsHeart_server/utils/rand_str"
)

func HandleQueryUnregAgents(c *gin.Context) {
	all, err := agent.GetAllUnreg()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": all,
	})
}

func HandleStartAgents(c *gin.Context) {
	var a agent.Agent
	err := c.ShouldBindJSON(&a)
	if err != nil {
		// client ip: who request this url
		logger.RigisterLog.Errorf("action=start agent;do=bind json;client_ip=%s;err=%s", c.ClientIP(), err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if a.Status == agent.ACCEPTED {
		a.Token = rand_str.GetStr(20)
	}

	// save to db
	if err = a.ChangeStatus(); err != nil {
		// agent ip: which will be start up
		logger.RigisterLog.Errorf("action=start agent;do=save db;agent_ip=%s;err=%s", a.RemoteAddr, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// tell agent the status: accept or deny
	err = a.StartUpAgent()
	if err != nil {
		logger.RigisterLog.Errorf("action=start agent;do=post agent;agent_ip=%s;err=%s", a.RemoteAddr, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
