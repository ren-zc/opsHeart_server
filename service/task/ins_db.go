package task

import (
	"errors"
	"opsHeart_server/db"
)

func (ins *TaskInstance) Create() error {
	return db.DB.Model(&TaskInstance{}).Create(ins).Error
}

func (ins *TaskInstance) Update(fields []string, values ...interface{}) error {
	if len(fields) != len(values) {
		return errors.New("length of fields not equal to length of values")
	}
	m := make(map[string]interface{})
	for i, f := range fields {
		m[f] = values[i]
	}
	return db.DB.Model(ins).Updates(m).Error
}

func (ins *TaskInstance) UpdateAndCheckCallbackVGROUP(fields []string, values ...interface{}) (c bool, err error) {
	if len(fields) != len(values) {
		err = errors.New("length of fields not equal to length of values")
		return
	}
	m := make(map[string]interface{})
	for i, f := range fields {
		m[f] = values[i]
	}

	tx := db.DB.Begin()

	err = tx.Model(ins).Updates(m).Error
	if err != nil {
		return
	}

	var allBrother []TaskInstance
	err = tx.Model(ins).Where("parent_ins_id = ?", ins.ParentInsID).
		Find(&allBrother).Error
	if err != nil {
		return
	}

	tx.Commit()

	if len(allBrother) < ins.BrotherNum {
		return
	}

	last := 0
	for _, v := range allBrother {
		//fmt.Printf("continue: %d, status: %d\n", v.ContinueOnFail, v.Status)
		if v.Status < STAGEFAILED {
			last++
		}
	}

	c = last == 0
	return
}

func (ins *TaskInstance) GetAllAgents() (sa []TaskStageAgent, err error) {
	err = db.DB.Model(ins).Related(&sa, "TaskStageAgents").Error
	return
}

func (ins *TaskInstance) GetAllStagesOfTask(tID uint) (ai []TaskInstance, err error) {
	err = db.DB.Model(ins).Where("parent_ins_id = ? and task_id = ?",
		ins.ParentInsID, tID).Find(&ai).Error
	return
}

func (ins *TaskInstance) GetBrotherInsByTaskID(i uint) (all []TaskInstance, err error) {
	err = db.DB.Model(ins).Where("parent_ins_id = ? and task_id = ?", ins.ParentInsID, i).
		Find(&all).Error
	return
}

func (ins *TaskInstance) GetNextStageInstance() (ti []TaskInstance, err error) {
	err = db.DB.Model(ins).
		Where("name = ? and parent_ins_id = ? and task_id = ?",
			ins.Name, ins.ParentInsID, ins.TaskID).
		Find(&ti, "stage_seq = ?", ins.StageSeq+1).Error
	return
}

func (ins *TaskInstance) GetAllChildIns() (all []TaskInstance, err error) {
	err = db.DB.Model(ins).Where("parent_ins_id = ?", ins.ID).Find(&all).Error
	return
}

func GetAllInsByNameAndTaskID(n string, i uint) (allIns []TaskInstance, err error) {
	err = db.DB.Model(&TaskInstance{}).Where("name = ? and task_id = ?", n, i).
		Find(&allIns).Error
	return
}

func GetInstanceByID(id uint) (p TaskInstance, err error) {
	err = db.DB.Model(&TaskInstance{}).
		Where("id = ?", id).First(&p).Error
	return
}

func GetAllInsByName(name string) (all []TaskInstance, err error) {
	err = db.DB.Model(&TaskInstance{}).Find(&all, "name = ?", name).Error
	return
}
