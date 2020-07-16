package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"opsHeart_server/common"
	"opsHeart_server/logger"
	"opsHeart_server/service/task"
)

func HandleRunTask(c *gin.Context) {
	var tk task.Task
	err := c.ShouldBindJSON(&tk)
	if err != nil {
		logger.TaskLog.Errorf("action=handle front run task request;do=bind data;err=%s;err_code=%d",
			err.Error(), common.BindPostDataErr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    err.Error(),
			"err_code": common.BindPostDataErr,
		})
		return
	}
	tkInDb, _ := task.GetTaskByID(tk.ID)

	// the customized target hosts
	if tk.CollectionType != 0 && tk.CollectionValue != "" {
		tkInDb.CollectionType = tk.CollectionType
		tkInDb.CollectionValue = tk.CollectionValue
	}

	insName := task.NewInsName(tk.ID)
	parentIns := task.TaskInstance{
		Name: insName,
	}

	err = tkInDb.Run(&parentIns)
	if err != nil {
		logger.TaskLog.Errorf("action=handle front run task;do=run task;err=%s;err_code=%d",
			err.Error(), common.RunTaskErr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    err.Error(),
			"err_code": common.RunTaskErr,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"err_code": common.NoError,
		"status":   "success",
	})
}
