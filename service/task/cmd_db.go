package task

import "opsHeart/db"

func (c *TaskCmd) QueryByTaskID() (err error) {
	err = db.DB.Model(c).Where("task_id = ?", c.TaskID).First(c).Error
	return
}

func (c *TaskCmd) Create() (err error) {
	err = db.DB.Create(c).Error
	return
}
