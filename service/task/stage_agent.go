package task

import "opsHeart_server/db"

func (sa *TaskStageAgent) Update(f string, v interface{}) error {
	return db.DB.Model(&TaskStageAgent{}).
		Where("id = ?", sa.ID).Update(f, v).Error
}

func (sa *TaskStageAgent) Create() error {
	return db.DB.Model(&TaskStageAgent{}).Create(sa).Error
}
