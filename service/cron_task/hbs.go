package cron_task

import (
	"opsHeart/db"
	"opsHeart/logger"
	"opsHeart/service/agent"
	"opsHeart/utils/cron"
	"time"
)

var CheckHbs *cron.Cr

func FindOffinedAgent() error {
	t := time.Now().Add(-3 * time.Minute)
	return db.DB.Model(&agent.Agent{}).
		Where("hbs_time < ? AND hbs_status = ?", t, agent.WORKING).
		Update("hbs_status", agent.OFFLINE).Error
}

func init() {
	var err error
	CheckHbs, err = cron.NewCron(FindOffinedAgent, nil, 3*time.Minute, 1)
	if err != nil {
		logger.HbsLog.Errorf("action=start check hbs cron;err=%s", err.Error())
	}
}
