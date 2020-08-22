package task

import (
	"opsHeart_server/logger"
	"time"
)

func VgroupChildChecker(ins *TaskInstance, tk *Task) {
	for {
		time.Sleep(2 * time.Second)
		allChild, err := ins.GetAllChildInsWithUnfinished()
		if err != nil {
			logger.TaskLog.Errorf("action=cron check vgroup;ins_id=%d;err=%s",
				ins.ID, err.Error())
		}

		// unfinished children length
		n := len(allChild)

		if n > 0 {
			continue
		}

		if n == 0 {
			ins.StageFinish(tk)
			logger.TaskLog.Infof("action=vgroup end;ins_id=%d;task_id=%d", ins.ID, tk.ID)
			break
		}
	}
}
