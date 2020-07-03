package task

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"opsHeart/logger"
	"opsHeart/service/collection"
	"opsHeart/utils/rand_str"
	"strconv"
	"time"
)

type StageAgentsMap struct {
	stage  int
	agents []string
}

// NewInsName create a instance name by task id
func NewInsName(taskID uint) string {
	randStr := rand_str.GetStr(4)
	tm := time.Now()
	InsName := fmt.Sprintf("%d-%d%d%d%d%d%d-%s",
		taskID, tm.Year(), tm.Month(), tm.Day(),
		tm.Hour(), tm.Minute(), tm.Second(), randStr)
	return InsName
}

// NewStageName create a stage name by task id
func NewStageName(taskID uint, stg int) string {
	return fmt.Sprintf("%d-%d-%s", taskID, stg, rand_str.GetStr(4))
}

// getCollection get all target agents IPs by collection definition.
func (t *Task) getCollection(parentIns *TaskInstance) ([]string, error) {
	ct := t.CollectionType
	cv := t.CollectionValue

	switch ct {
	case CollList:
		var allIPs []string
		err := json.Unmarshal([]byte(cv), &allIPs)
		if err != nil {
			return nil, err
		}
		return allIPs, nil
	case CollName:
		// get collection by agent collection rule
		ac, err := collection.QueryCollByName(cv)
		if err != nil {
			return nil, err
		}
		return ac.GetAllIPs()
	case CollInheritPercent:
		n, err := strconv.Atoi(cv)
		if err != nil {
			return nil, err
		}
		return parentIns.GetAllInsIPsByPercentOrNum(n, StagePercent, t.SplitParent)
	case CollInheritNum:
		n, err := strconv.Atoi(cv)
		if err != nil {
			return nil, err
		}
		return parentIns.GetAllInsIPsByPercentOrNum(n, StageNumber, t.SplitParent)
	}
	return nil, nil
}

// SplitIPsList split IPs list to several parts by stage list.
// return map[stage_num][]IP
func SplitIPsList(st StType, sNum []int, ips []string) (sam []StageAgentsMap, err error) {
	ipsL := len(ips)
	sNumL := len(sNum)
	if ipsL == 0 {
		return nil, errors.New("ip list 0 length")
	}
	if sNumL == 0 {
		return nil, errors.New("stage list 0 length")
	}

	sNewNum := make([]int, len(sNum))

	if st == StagePercent {
		for i, v := range sNum {
			f := float64(v) / 100
			n := int(math.Ceil(f * float64(len(ips))))
			if n >= ipsL {
				sNewNum[i] = ipsL
				break
			}
			sNewNum[i] = n
			ipsL -= n
		}
	}

	if st == StageNumber {
		for i, v := range sNum {
			if v >= ipsL {
				sNewNum[i] = ipsL
				break
			}
			sNewNum[i] = v
			ipsL -= v
		}
	}

	startOffset := 0
	for i, v := range sNewNum {
		sa := StageAgentsMap{
			stage:  sNum[i],
			agents: ips[startOffset : startOffset+v],
		}
		sam = append(sam, sa)
		startOffset += v
	}

	return
}

// Run start a task
// before Run, should run NewInsName to get a instance name
func (t *Task) Run(parentIns *TaskInstance) error {
	st := t.StageType
	sg := t.Stages
	var stages []int
	err := json.Unmarshal([]byte(sg), &stages)
	if err != nil {
		return err
	}

	// get ips
	IPs, err := t.getCollection(parentIns)
	if err != nil {
		return err
	}

	// split IP list by stage list
	samList, err := SplitIPsList(st, stages, IPs)
	if err != nil {
		return err
	}

	// create stage name and save to stage agent data db
	var firstStageIns TaskInstance
	var stageSeqNum uint
	for _, v := range samList {
		// save stage agent ip to db
		stageRef := NewStageName(t.ID, v.stage)
		IPList := v.agents
		for _, ip := range IPList {
			sa := TaskStageAgent{
				StageName: stageRef,
				IP:        ip,
			}
			err := sa.Create()
			if err != nil {
				return err
			}
		}

		// save task stage ins data to db
		stageSeqNum++
		ti := TaskInstance{
			Name:           parentIns.Name,
			ParentInsID:    parentIns.ID,
			TaskID:         t.ID,
			Stage:          v.stage,
			StageSeq:       stageSeqNum,
			StageAgents:    stageRef,
			RunBy:          "reserved",
			ParentIsX:      t.ParentXGROUP,
			ParentIsV:      t.ParentVGROUP,
			BrotherNum:     t.ChildesNum,
			ContinueOnFail: t.ContinueOnFail,
			//IsVGROUP:    parentIns.IsVGROUP,
		}

		err := ti.Create()
		if err != nil {
			return err
		}

		if stageSeqNum == 1 {
			firstStageIns = ti
		}
	}

	// start first stage
	logger.TaskLog.Debugf("*** task run: %d, first stage: %d\n", t.ID, firstStageIns.ID)

	go firstStageIns.StartStage(t)

	logger.TaskLog.Infof("action=run task;task_name=%s;task_id=%d;ins_name=%s;",
		t.Name, t.ID, parentIns.Name)

	return nil
}
