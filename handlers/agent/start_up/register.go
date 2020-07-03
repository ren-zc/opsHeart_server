package start_up

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opsHeart/logger"
	"opsHeart/service/agent"
)

// HandleAgentRegister used for agent register self to server
func HandleAgentRegister(c *gin.Context) {
	ag := agent.Agent{}
	ag.Status = agent.REGISTER

	var agentIP string
	agentIP = c.ClientIP()

	err := c.ShouldBindJSON(&ag)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		logger.RigisterLog.Error("get register data err: %s, from: %s", err.Error(), agentIP)
		return
	}

	ag.RemoteAddr = agentIP

	//var tmp v1.Agent
	//err = db.DB.FirstOrCreate(&tmp, v1.Agent{UUID: ag.UUID}).Error
	if ag.IsExist() {
		err = ag.UpdateDat()
	} else {
		err = ag.InsertDat()
	}
	if err != nil {
		logger.ServerLog.Errorf("uuid=%s, hostname=%s, addr=%s, register_error=%s",
			ag.UUID, ag.Hostname, ag.RemoteAddr, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal error",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success", // mean: got register request successfully.
	})

	logger.RigisterLog.Infof("action=register agent;uuid=%s, hostname=%s, addr=%s, status=success",
		ag.UUID, ag.Hostname, ag.RemoteAddr)
}
