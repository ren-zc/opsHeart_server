package task

import "opsHeart_server/db"

func (s *TaskScript) QueryByTaskID() (err error) {
	err = db.DB.Model(s).Where("task_id = ?", s.TaskID).First(s).Error
	return
}

func (s *TaskScript) Create() (err error) {
	err = db.DB.Create(s).Error
	return
}
