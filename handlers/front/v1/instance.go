package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opsHeart_server/common"
	"opsHeart_server/logger"
	"opsHeart_server/service/task"
)

func HandleRunIns(c *gin.Context) {
	var tiTmp task.TaskInstance
	err := c.ShouldBindJSON(&tiTmp)
	if err != nil {
		logger.TaskLog.Errorf("action=handle run instance;do=bind tk data;err=%s;err_code=%d",
			err.Error(), common.BindPostDataErr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    err.Error(),
			"err_code": common.BindPostDataErr,
		})
		return
	}

	tk, err := task.GetTaskByID(tiTmp.TaskID)
	if err != nil {
		logger.TaskLog.Errorf("action=handle run instance;do=get task;err=%s;err_code=%d",
			err.Error(), common.QueryTaskByIDErr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    err.Error(),
			"err_code": common.QueryTaskByIDErr,
		})
		return
	}

	ti, err := task.GetInstanceByID(tiTmp.ID)
	if err != nil {
		logger.TaskLog.Errorf("action=handle run instance;do=get instance;err=%s;err_code=%d",
			err.Error(), common.QueryInstanceByIDErr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    err.Error(),
			"err_code": common.QueryInstanceByIDErr,
		})
		return
	}

	go ti.StartStage(&tk)

	c.JSON(http.StatusOK, gin.H{
		"err_code": common.NoError,
		"status":   "success",
	})
}
