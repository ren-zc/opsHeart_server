package agent

import (
	"opsHeart/db"
	"time"
)

func (a *Agent) UpdateHbs() error {
	return db.DB.Model(&Agent{}).Where("uuid", a.UUID).
		Update("hbs_status", WORKING).
		Update("hbs_time", time.Now()).Error
}
