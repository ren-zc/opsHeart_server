package task

import (
	"errors"
	"fmt"
	"math"
	"opsHeart_server/common"
	"opsHeart_server/db"
	"opsHeart_server/logger"
	"time"
)

// GetAllInsIPs get instance id stage all related agent IPs
func (ins *TaskInstance) GetAllInsIPs() (allIPs []string, err error) {
	stageAgents, err := ins.GetAllAgents()
	for _, v := range stageAgents {
		allIPs = append(allIPs, v.IP)
	}
	return
}

// GetAllInsIPsByPercentOrNum get instance related agent IPs by percent or number
func (ins *TaskInstance) GetAllInsIPsByPercentOrNum(p int, t StType, cs splitColl) ([]string, error) {
	var stageAgents []TaskStageAgent
	var agents []string

	err := db.DB.Model(ins).Related(&stageAgents, "TaskStageAgents").Error
	if err != nil {
		return nil, err
	}

	var n int
	if t == StagePercent {
		n = int(math.Ceil((float64(p) / 100) * float64(len(stageAgents))))
	} else {
		n = p
	}

	for _, v := range stageAgents {
		if cs == DoSplit && v.ChildUse == IsUsed {
			continue
		}
		if n == 0 {
			break
		}
		agents = append(agents, v.IP)
		if cs == DoSplit {
			err := v.Update("child_use", IsUsed)
			if err != nil {
				return nil, err
			}
		}
		n--
	}
	return agents, nil
}

func (ins *TaskInstance) logFailMsg(status StageStatus, msg string, t *Task) {
	logger.TaskLog.Debugf(">>> ins: %d, update status in log fail: %d\n", ins.ID, status)

	_ = ins.Update([]string{"status", "ins_msg", "end_at"}, status, msg, time.Now())
	logger.TaskLog.Error(msg)

	ins.Status = status

	if ins.ContinueOnFail == 1 {
		ins.StageFinish(t)
	} else {
		ins.CallbackParentStage()
	}
}

func (ins *TaskInstance) runChildTask(tk *Task, seq uint) {
	c, err := tk.GetTheSeqChild(seq)
	if err != nil {
		errMsg := fmt.Sprintf("action=start stage;do=get seq %d child task;err=%s",
			seq, err.Error())
		ins.logFailMsg(STAGEFAILED, errMsg, tk)
		return
	}
	if len(c) != 1 {
		err = errors.New("multi child task exist")
		errMsg := fmt.Sprintf("action=start stage;do=get seq %d child task;err=%s",
			seq, err.Error())
		ins.logFailMsg(STAGEFAILED, errMsg, tk)
		return
	}
	err = c[0].Run(ins)
	if err != nil {
		errMsg := fmt.Sprintf("action=start stage;do=run seq %d child task;err=%s",
			seq, err.Error())
		ins.logFailMsg(STAGEFAILED, errMsg, tk)
	}
}

func (ins *TaskInstance) StartStage(t *Task) {
	tm := time.Now()
	logger.TaskLog.Debugf(">>> ins: %d, update status in start stage: %d\n",
		ins.ID, STAGERUNNING)
	err := ins.Update([]string{"status", "start_at"}, STAGERUNNING, tm)
	if err != nil {
		logger.TaskLog.
			Errorf("action=start stage;do=save stage status;ins_id=%d;ins_name=%s;task_id=%d;err=%s",
				ins.ID, ins.Name, ins.TaskID, err.Error())
		return
	}

	ins.Status = STAGERUNNING
	ins.StartAt = tm

	// check continue condition
	ok, err := ins.checkContinueCondition(t)
	if err != nil {
		errMsg := fmt.Sprintf("action=start stage;do=check continue condition;ins_id=%d;ins_name=%s;task_id=%d;err=%s",
			ins.ID, ins.Name, ins.TaskID, err.Error())
		ins.logFailMsg(STAGEFAILED, errMsg, t)
		return
	}

	if !ok {
		errMsg := fmt.Sprintf("action=start stage;do=check continue condition;ins_id=%d;ins_name=%s;task_id=%d;err=continue check fail",
			ins.ID, ins.Name, ins.TaskID)
		ins.logFailMsg(STAGEFAILED, errMsg, t)
		return
	}

	switch t.TkType {
	case TASKROOT:
		// get the first child task to run
		ins.runChildTask(t, 1)
		return
	case HGROUP:
		// get the first child task to run
		ins.runChildTask(t, 1)
		return
	case VGROUP:
		// get all child task to run
		allChilds, err := t.GetAllChildTask()
		if err != nil {
			errMsg := fmt.Sprintf("action=start stage;do=get child tasks;ins_id=%d;ins_name=%s;err=%s",
				ins.ID, ins.Name, err.Error())
			ins.logFailMsg(STAGEFAILED, errMsg, t)
			return
		}

		cl := len(allChilds)
		for _, c := range allChilds {
			tk := c
			tk.ParentVGROUP = 1
			tk.ChildesNum = cl
			err := tk.Run(ins)
			if err != nil {
				errMsg := fmt.Sprintf("action=start stage;do=start child task;ins_id=%d;ins_name=%s;child_task_id=%d;task_id=%s;err=%s",
					ins.ID, ins.Name, tk.ID, t.Name, err.Error())
				ins.logFailMsg(STAGEFAILED, errMsg, t)
				return
			}
		}

		// start cron
		VgroupChildChecker(ins, t)

	case XGROUP:
		// get the previous task result and determine which child task to run
		rst, err := ins.getPreviousInsRst(t)
		if err != nil {
			errMsg := fmt.Sprintf("action=start stage;do=get brother task stages;ins_id=%d;ins_name=%s",
				ins.ID, ins.Name)
			ins.logFailMsg(STAGEFAILED, errMsg, t)
			return
		}

		// get all child task
		childes, err := t.GetAllChildTask()
		if err != nil {
			errMsg := fmt.Sprintf("action=start stage;do=xgroup get all childs;stage=%d;ins_id=%d;ins_name=%s",
				ins.Stage, ins.ID, ins.Name)
			ins.logFailMsg(STAGEFAILED, errMsg, t)
			return
		}

		if len(childes) != 2 {
			errMsg := fmt.Sprintf("action=start stage;do=check xgroup childs length;stage=%d;ins_id=%d;ins_name=%s",
				ins.Stage, ins.ID, ins.Name)
			ins.logFailMsg(STAGEFAILED, errMsg, t)
			return
		}

		var tSuccess Task
		var tFailed Task
		for _, v := range childes {
			if v.SeqNum == 1 {
				tSuccess = v
				tSuccess.ParentXGROUP = 1
			}
			if v.SeqNum == 2 {
				tFailed = v
				tFailed.ParentXGROUP = 1
			}
		}

		// choose one to run
		if rst == STAGESUCCESS {
			err := tSuccess.Run(ins)
			if err != nil {
				errMsg := fmt.Sprintf("action=start stage;task_id=%d;task_name=%s;do=run 1st task;ins_id=%d;ins_name=%s;err=%s",
					tSuccess.ID, tSuccess.Name, ins.ID, ins.Name, err.Error())
				ins.logFailMsg(STAGEFAILED, errMsg, t)
			}
			return
		}
		if rst == STAGEFAILED {
			err := tFailed.Run(ins)
			if err != nil {
				errMsg := fmt.Sprintf("action=start stage;task_id=%d;task_name=%s;do=run 2nd task;ins_id=%d;ins_name=%s;err=%s",
					tSuccess.ID, tSuccess.Name, ins.ID, ins.Name, err.Error())
				ins.logFailMsg(STAGEFAILED, errMsg, t)
				return
			}
		}

	default:
		ins.startBasicTask(t)
	}
}

func (ins *TaskInstance) startBasicTask(tk *Task) {
	var bt basicTask
	args, _ := ins.GetAllArgs()
	switch tk.TkType {
	case TASKCMD:
		tc := &TaskCmd{}
		tc.TaskID = ins.TaskID
		tc.Args = args
		bt = tc
	case TASKSCRIPT:
		ts := &TaskScript{}
		ts.TaskID = ins.TaskID
		ts.Args = args
		bt = ts
	case SYNCFILE:
		tsf := &TaskSyncFile{}
		tsf.TaskID = ins.TaskID
		bt = tsf
	default:
		return
	}

	// start task and block until task finish
	err := bt.QueryByTaskID()
	if err != nil {
		errMsg := fmt.Sprintf("action=start stage;do=query cmd task;ins_id=%d;ins_name=%s;task_id=%d;err=%s",
			ins.ID, ins.Name, ins.TaskID, err.Error())
		ins.logFailMsg(STAGEFAILED, errMsg, tk)
		return
	}
	agents, _ := ins.GetAllInsIPs()
	err = bt.start(&agents)
	if err != nil {
		errMsg := fmt.Sprintf("action=start stage;do=start cmd task;ins_id=%d;ins_name=%s;task_id=%d;err=%s",
			ins.ID, ins.Name, ins.TaskID, err.Error())
		ins.logFailMsg(STAGEFAILED, errMsg, tk)
		return
	}

	ins.StageFinish(tk)
}

// getPreviousInsRst get the status of previous task
func (ins *TaskInstance) getPreviousInsRst(t *Task) (StageStatus, error) {
	allStages, err := ins.getAllStagesOfBrotherTask(t)
	if err != nil {
		return 0, err
	}

	rst := STAGESUCCESS
	for _, v := range allStages {
		if v.Status != STAGESUCCESS {
			rst = STAGEFAILED
			break
		}
	}
	return rst, nil
}

func (ins *TaskInstance) checkContinueCondition(t *Task) (bool, error) {
	// no continue task be set
	if t.ContinueByTask == 0 {
		return true, nil
	}

	allBrothers, err := ins.GetBrotherInsByTaskID(t.ContinueByTask)
	if err != nil {
		return false, err
	}
	for _, v := range allBrothers {
		if v.Status != t.ContinueRst {
			return false, errors.New("mismatched continue result")
		}
	}

	return true, nil
}

func (ins *TaskInstance) getAllStagesOfBrotherTask(t *Task) (ai []TaskInstance, err error) {
	tb, err := t.GetOldBrotherTask()
	if err != nil {
		return
	}
	return ins.GetAllStagesOfTask(tb.ID)
}

// StageFinish: when stage finished, it will be called
func (ins *TaskInstance) StageFinish(t *Task) {
	childIns, _ := ins.GetAllChildIns()

	// get the status
	insStatus := STAGESUCCESS
	for _, v := range childIns {
		if v.Status == STAGEFAILED {
			insStatus = STAGEFAILED
			break
		}
	}

	// save status to db
	if ins.Status < STAGEFAILED {
		var err error
		//var callback bool
		logger.TaskLog.Debugf(">>> ins: %d, update status in finish stage: %d\n",
			ins.ID, insStatus)
		err = ins.Update([]string{"status", "end_at"}, insStatus, time.Now())
		if err != nil {
			logger.TaskLog.Errorf("action=start stage;do=update success;ins_id=%d;ins_name=%s;err=%s",
				ins.ID, ins.Name, err.Error())
			return
		}
	}

	if insStatus == STAGEFAILED && ins.ContinueOnFail != 1 {
		ins.CallbackParentStage()
		return
	}

	// start next stage
	if !common.StepPause {
		runNextIns := ins.RunNextStage(t)
		if !runNextIns {
			return
		}
	}

	// if task in a vgroup task, return
	if ins.ParentIsV == 1 {
		return
	}

	// if it's a task root, exit, because all stage finished.
	if t.TkType == TASKROOT {
		logger.TaskLog.Debugf("*** ins: %d, task id: %d, task is root, return.\n",
			ins.ID, t.ID)
		return
	}

	// if no next stage, call next task
	parentIns, _ := GetInstanceByID(ins.ParentInsID)

	// if parent instance is a xgroup task, it can only call finish to callback parent
	//if ins.ParentIsX != 1 {
	//}
	continueOn := ins.RunNextTask(t, &parentIns)
	if !continueOn {
		return
	}

	logger.TaskLog.Debugf("*** ins: %d, callback parent 1: %d\n", ins.ID, parentIns.ID)

	// call parent instance stage finish
	parentTk, _ := GetTaskByID(parentIns.TaskID)
	parentIns.StageFinish(&parentTk)
	return
}

func (ins *TaskInstance) RunNextTask(t *Task, parentIns *TaskInstance) (continueOn bool) {
	nextTk, _ := t.GetNextTask()
	nextTkL := len(nextTk)
	if nextTkL == 1 {
		logger.TaskLog.Debugf("*** ins: %d, next task: %d\n", ins.ID, nextTk[0].ID)
		err := nextTk[0].Run(parentIns)
		if err != nil {
			logger.TaskLog.Errorf("action=stage finish;do=start task;task_id=%d;task_name=%s;ins_id=%d;ins_name=%s;err=%s;",
				nextTk[0].ID, nextTk[0].Name, ins.ID, ins.Name, err.Error())
			return false
		}
		return false
	}
	if nextTkL > 1 {
		logger.TaskLog.Errorf("action=stage finish;do=start task;ins_id=%d;ins_name=%s;err=more task found",
			ins.ID, ins.Name)
		return false
	}
	return true
}

func (ins *TaskInstance) RunNextStage(t *Task) (nextTask bool) {
	nextIns, _ := ins.GetNextStageInstance()
	nextInsL := len(nextIns)
	fmt.Printf("next ins length: %d\n", nextInsL)
	if nextInsL == 1 {
		if nextIns[0].Status >= STAGERUNNING {
			logger.TaskLog.Debugf("*** ins: %d, next ins: %d, next ins is running, return\n",
				ins.ID, nextIns[0].ID)
			return
		}
		logger.TaskLog.Debugf("*** ins: %d, next ins: %d\n", ins.ID, nextIns[0].ID)

		// stage not auto start but step pause is false set by config file.
		nextIns[0].StartStage(t)
		logger.TaskLog.Infof("action=run next stage;task_name=%s;task_id=%d;ins=%d;next_ins=%d",
			t.Name, t.ID, ins.ID, nextIns[0].ID)

		return false
	}
	if nextInsL > 1 {
		logger.TaskLog.Errorf("action=stage finish;ins_id=%d;ins_name=%s;err=there are %d stage for seq %d",
			ins.ID, ins.Name, nextInsL, ins.StageSeq+1)
		return false
	}
	return true
}

func (ins *TaskInstance) CallbackParentStage() {
	parentIns, _ := GetInstanceByID(ins.ParentInsID)
	parentTk, _ := GetTaskByID(parentIns.TaskID)
	logger.TaskLog.Debugf("*** ins: %d, callback parent 0: %d\n", ins.ID, parentIns.ID)
	if parentIns.ID == 0 {
		return
	}
	parentIns.StageFinish(&parentTk)
}
