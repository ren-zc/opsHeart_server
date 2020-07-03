package task

import (
	"github.com/pkg/errors"
	"opsHeart/db"
)

// GetTheSeqChild get the seq-th child of task t
func (t *Task) GetTheSeqChild(seq uint) ([]Task, error) {
	var childs []Task
	err := db.DB.Model(t).Where("parent_task_id = ?", t.ID).
		Find(&childs, "seq_num = ?", seq).Error
	return childs, err
}

func GetTaskByID(i uint) (t Task, err error) {
	db.DB.Model(t).First(&t, "id = ?", i)
	return
}

func (t *Task) GetOldBrotherTask() (bt Task, err error) {
	if t.SeqNum == 1 {
		err = errors.New("the seq 1 invalid for finding previous brother task")
		return
	}
	err = db.DB.Model(t).Where("parent_task_id = ?", t.ParentTaskID).
		First(&bt, "seq_num = ?", t.SeqNum-1).Error
	return
}

func (t *Task) GetAllChildTask() (allChild []Task, err error) {
	err = db.DB.Model(t).Where("parent_task_id = ?", t.ID).
		Find(&allChild).Error
	return
}

func (t *Task) CheckRootNameIsExist() bool {
	var tmp Task
	db.DB.Model(t).Where("name = ? and tk_type = ?", t.Name, 0).First(&tmp)
	return tmp.Name == t.Name
}

// Before create, if continue depending is set, should check the depended task is a brother node.
func (t *Task) Create() error {
	if t.TkType == TASKROOT {
		if t.CheckRootNameIsExist() {
			return errors.New("root task has exist")
		}
	}
	return db.DB.Create(t).Error
}

func (t *Task) GetNextTask() (next []Task, err error) {
	err = db.DB.Model(t).Where("parent_task_id = ?", t.ParentTaskID).
		Find(&next, "seq_num = ?", t.SeqNum+1).Error
	return
}

// CheckNameIsExist check root task name is unique
// for `root` task, task name must be unique
func CheckNameIsExist(n string) (ok bool) {
	t := Task{}
	err := db.DB.Model(&t).Where("tk_type = 0").First(&t, "name = ?", n).Error
	if err != nil {
		return false
	}
	if t.Name == n {
		return true
	}
	return false
}
