package task

import "opsHeart_server/db"

func (sf *TaskSyncFile) QueryByTaskID() (err error) {
	err = db.DB.Model(sf).Where("task_id = ?", sf.TaskID).First(sf).Error
	return
}

func (sf *TaskSyncFile) Create() (err error) {
	err = db.DB.Create(sf).Error
	return
}
