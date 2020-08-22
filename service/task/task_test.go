package task

import (
	"encoding/json"
	"opsHeart_server/common"
	"opsHeart_server/conf"
	"opsHeart_server/db"
	"opsHeart_server/logger"
	"opsHeart_server/utils/rand_str"
	"strconv"
	"testing"
	"time"
)

func TestSplitIPSList(t *testing.T) {
	ipsL := 83
	ips := make([]string, 0)
	for i := 0; i < ipsL; i++ {
		rs := rand_str.GetStr(4)
		ips = append(ips, rs)
	}
	t.Logf("ips len: %d\n", len(ips))
	t.Logf("ips: %v\n", ips)

	type SplitTest struct {
		t StType
		s []int
	}

	t.Log("=======t1 number======")
	// test 1
	t1 := SplitTest{
		t: StageNumber,
		s: []int{10, 20, 30, 50},
	}
	tm1, err := SplitIPsList(t1.t, t1.s, ips)
	if err != nil {
		t.Fatalf("t1 test: %s\n", err)
	}
	//t.Logf("t1 test: %v\n", tm1)
	for _, v := range tm1 {
		t.Log(v.stage)
		t.Logf("len: %d, list: %v", len(v.agents), v.agents)
	}

	t.Log("=======t2 number======")
	// test2
	t2 := SplitTest{
		t: StageNumber,
		s: []int{10, 20, 9},
	}
	tm2, err := SplitIPsList(t2.t, t2.s, ips)
	if err != nil {
		t.Fatalf("t2 test: %s\n", err)
	}
	//t.Logf("t2 test: %v\n", tm2)
	for _, v := range tm2 {
		t.Log(v.stage)
		t.Logf("len: %d, list: %v", len(v.agents), v.agents)
	}

	t.Log("=======t3 percent======")
	// test3
	t3 := SplitTest{
		t: StagePercent,
		s: []int{10, 30, 40, 50},
	}
	tm3, err := SplitIPsList(t3.t, t3.s, ips)
	if err != nil {
		t.Fatalf("t3 test: %s\n", err)
	}
	//t.Logf("t3 test: %v\n", tm3)
	for _, v := range tm3 {
		t.Log(v.stage)
		t.Logf("len: %d, list: %v", len(v.agents), v.agents)
	}

	t.Log("=======t4 percent======")
	// test4
	t4 := SplitTest{
		t: StagePercent,
		s: []int{10, 30, 20, 5, 8},
	}
	tm4, err := SplitIPsList(t4.t, t4.s, ips)
	if err != nil {
		t.Fatalf("t3 test: %s\n", err)
	}
	//t.Logf("t4 test: %v\n", tm4)
	for _, v := range tm4 {
		t.Log(v.stage)
		t.Logf("len: %d, list: %v", len(v.agents), v.agents)
	}

	t.Log("=======t5 default percent======")
	// test5
	t5 := SplitTest{
		t: StagePercent,
		s: []int{100},
	}
	tm5, err := SplitIPsList(t5.t, t5.s, ips)
	if err != nil {
		t.Fatalf("t5 test: %s\n", err)
	}
	//t.Logf("t5 test: %v\n", tm5)
	for _, v := range tm5 {
		t.Log(v.stage)
		t.Logf("len: %d, list: %v", len(v.agents), v.agents)
	}
}

func TestTask_Run(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	common.InitRunningEnv()
	db.InitDB()
	db.DB.AutoMigrate(&Task{})
	db.DB.AutoMigrate(&TaskInstance{})
	db.DB.AutoMigrate(&TaskCmd{})
	db.DB.AutoMigrate(&TaskScript{})
	db.DB.AutoMigrate(&TaskStageAgent{})
	db.DB.AutoMigrate(&TaskArg{})
	db.DB.AutoMigrate(&InsArg{})

	ipList := make([]string, 0)
	for i := 0; i < 20; i++ {
		rands := rand_str.GetStr(4)
		ipList = append(ipList, rands)
	}

	collValByte, _ := json.Marshal(ipList)
	stageInt := []int{100}
	sg, _ := json.Marshal(stageInt)

	trStage := []int{50, 50}
	trStageByte, _ := json.Marshal(trStage)

	// *** ROOT TASK ***
	tr := Task{
		Name:            "test task 93",
		TkType:          TASKROOT,
		CollectionType:  CollList,
		CollectionValue: string(collValByte),
		Stages:          string(trStageByte), // 50% 50%
		//Stages:          string(sg), // 100%
		Desc: "a basic test task",
		//ContinueOnFail: 1,
	}
	err = tr.Create()
	if err != nil {
		t.Fatalf("test create tr err: %s", err.Error())
	}

	// *** THE FIRST TASK ***
	t1Stages := []int{50, 50}
	stagesByte, _ := json.Marshal(t1Stages)

	t1 := Task{
		Name:            tr.Name,
		TkType:          TASKCMD,
		ParentTaskID:    tr.ID,
		SeqNum:          1,
		CollectionValue: strconv.Itoa(50),
		Stages:          string(stagesByte),
		SplitParent:     DoSplit,
		ContinueOnFail:  1,
		TaskArgs: []TaskArg{
			{
				TaskName: tr.Name,
				ArgName:  "t1Var1",
				ArgType:  COMMONSTR,
				ArgValue: "t1Var1Value",
			},
			{
				TaskName: tr.Name,
				ArgName:  "t1Var2",
				ArgType:  AGENTFACT,
				ArgValue: "t1Var2Value",
			},
			{
				TaskName: tr.Name,
				ArgName:  "t1Var3",
				ArgType:  AGENTTAG,
				ArgValue: "t1Var3Value",
			},
		},
	}
	err = t1.Create()
	if err != nil {
		t.Fatalf("test create t1 err: %s", err.Error())
	}
	c1 := TaskCmd{
		TaskID:  t1.ID,
		Cmd:     "echo1",
		Opt:     "123 ok",
		Timeout: 60,
	}
	err = c1.Create()
	if err != nil {
		t.Fatalf("test create c1 err: %s", err.Error())
	}

	// *** THE SECOND TASK ***
	t2Stage := []int{30, 70}
	t2stageByte, _ := json.Marshal(t2Stage)

	t2 := Task{
		//Name:           "tmp-not-used",
		Name:   tr.Name,
		TkType: VGROUP,
		//ParentTaskID:   1000000,
		ParentTaskID: tr.ID,
		SeqNum:       2,
		Stages:       string(t2stageByte),
		SplitParent:  DoNotSplit,
		//ContinueOnFail: 1,
	}
	err = t2.Create()
	if err != nil {
		t.Fatalf("test create t2 err: %s", err.Error())
	}

	t21 := Task{
		Name:            tr.Name,
		TkType:          TASKCMD,
		ParentTaskID:    t2.ID,
		SeqNum:          1,
		CollectionValue: strconv.Itoa(50),
		Stages:          string(sg),
		SplitParent:     DoSplit,
	}
	err = t21.Create()
	if err != nil {
		t.Fatalf("test create t21 err: %s", err.Error())
	}
	c21 := TaskCmd{
		TaskID: t21.ID,
		Cmd:    "test21",
		//Opt:     "123 ok test",
		Timeout: 60,
	}
	err = c21.Create()
	if err != nil {
		t.Fatalf("test create c21 err: %s", err.Error())
	}

	t22 := Task{
		Name:            tr.Name,
		TkType:          TASKSCRIPT,
		ParentTaskID:    t2.ID,
		SeqNum:          1,
		CollectionValue: strconv.Itoa(50),
		Stages:          string(sg),
		SplitParent:     DoSplit,
		TaskArgs: []TaskArg{
			{
				TaskName: tr.Name,
				ArgName:  "t22Var1",
				ArgType:  COMMONSTR,
				ArgValue: "t22Var1Value",
			},
			{
				TaskName: tr.Name,
				ArgName:  "t22Var2",
				ArgType:  AGENTFACT,
				ArgValue: "t22Var2Value",
			},
			{
				TaskName: tr.Name,
				ArgName:  "t22Var3",
				ArgType:  AGENTTAG,
				ArgValue: "t22Var3Value",
			},
		},
	}
	err = t22.Create()
	if err != nil {
		t.Fatalf("test create t22 err: %s", err.Error())
	}
	s22 := TaskScript{
		TaskID:  t22.ID,
		Shell:   "/bin/t22sh",
		Name:    "/tmp/test123.sh",
		RunAs:   "jacenr",
		Timeout: 60,
	}
	err = s22.Create()
	if err != nil {
		t.Fatalf("test create s22 err: %s", err.Error())
	}

	// *** THE THIRD TASK ***
	stageInt6040 := []int{40, 60}
	sg6040, _ := json.Marshal(stageInt6040)

	t3 := Task{
		Name:         tr.Name,
		TkType:       HGROUP,
		ParentTaskID: tr.ID,
		SeqNum:       3,
		//SeqNum:          2,
		CollectionValue: strconv.Itoa(50),
		Stages:          string(sg6040), // 40% 60%
		SplitParent:     DoSplit,
		ContinueOnFail:  1,
	}
	err = t3.Create()
	if err != nil {
		t.Fatalf("test create t3 err: %s", err.Error())
	}

	t31 := Task{
		Name:           tr.Name,
		TkType:         TASKCMD,
		ParentTaskID:   t3.ID,
		SeqNum:         1,
		Stages:         string(sg), // 100%
		SplitParent:    DoNotSplit,
		ContinueOnFail: 1,
		TaskArgs: []TaskArg{
			{
				TaskName: tr.Name,
				ArgName:  "t31Var1",
				ArgType:  COMMONSTR,
				ArgValue: "t31Var1Value",
			},
			{
				TaskName: tr.Name,
				ArgName:  "t31Var2",
				ArgType:  AGENTFACT,
				ArgValue: "t31Var2Value",
			},
			{
				TaskName: tr.Name,
				ArgName:  "t31Var3",
				ArgType:  AGENTTAG,
				ArgValue: "t31Var3Value",
			},
		},
	}
	err = t31.Create()
	if err != nil {
		t.Fatalf("test create t31 err: %s", err.Error())
	}
	c31 := TaskCmd{
		TaskID:  t31.ID,
		Cmd:     "echo31",
		Opt:     "123 ok test",
		Timeout: 60,
	}
	err = c31.Create()
	if err != nil {
		t.Fatalf("test create c31 err: %s", err.Error())
	}

	t32 := Task{
		Name:         tr.Name,
		TkType:       TASKSCRIPT,
		ParentTaskID: t3.ID,
		SeqNum:       2,
		Stages:       string(sg), // 100%
		SplitParent:  DoNotSplit,
	}
	err = t32.Create()
	if err != nil {
		t.Fatalf("test create t32 err: %s", err.Error())
	}
	s32 := TaskScript{
		TaskID:  t32.ID,
		Shell:   "/bin/csh",
		Name:    "/tmp/test123.sh",
		RunAs:   "jacenr",
		Timeout: 60,
	}
	err = s32.Create()
	if err != nil {
		t.Fatalf("test create s32 err: %s", err.Error())
	}

	t33 := Task{
		Name:           tr.Name,
		TkType:         TASKSCRIPT,
		ParentTaskID:   t3.ID,
		SeqNum:         3,
		Stages:         string(sg), // 100%
		SplitParent:    DoNotSplit,
		ContinueByTask: t31.ID,
		ContinueRst:    STAGEFAILED,
		ContinueOnFail: 1,
	}
	err = t33.Create()
	if err != nil {
		t.Fatalf("test create t33 err: %s", err.Error())
	}
	s33 := TaskScript{
		TaskID:  t33.ID,
		Shell:   "/bin/bash",
		Name:    "/tmp/test.sh",
		RunAs:   "jacenr",
		Timeout: 60,
	}
	err = s33.Create()
	if err != nil {
		t.Fatalf("test create s33 err: %s", err.Error())
	}

	t.Log("All tasks are created.")

	insName := NewInsName(tr.ID)
	ins := TaskInstance{
		Name: insName,
	}

	argMap := make(map[string]TaskArg)
	err = tr.InitInsArgs(argMap, insName)
	if err != nil {
		t.Fatalf("fail: %s", err)
	}

	err = tr.Run(&ins)
	if err != nil {
		t.Fatalf("test run tr fail: %s", err.Error())
	}

	allIns, err := GetAllInsByName(tr.Name)
	if err != nil {
		t.Fatalf("test run task fail: %s", err)
	}
	for _, v := range allIns {
		t.Logf("id:%d, parent_id:%d, task_id:%d, stage:%d, stage_seq:%d, start_at:%v, end_at:%v, status: %d, msg:%s\n",
			v.ID, v.ParentInsID, v.TaskID, v.Stage, v.StageSeq, v.StartAt, v.EndAt, v.Status, v.InsMsg)
	}

	time.Sleep(10 * time.Second)
}

func TestTask_Run2(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&Task{})
	db.DB.AutoMigrate(&TaskInstance{})
	db.DB.AutoMigrate(&TaskCmd{})
	db.DB.AutoMigrate(&TaskScript{})
	db.DB.AutoMigrate(&TaskStageAgent{})

	tr, _ := GetTaskByID(325)
	insName := NewInsName(325)
	parentIns := TaskInstance{
		Name: insName,
	}
	err = tr.Run(&parentIns)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}
	time.Sleep(10 * time.Second)
}

func TestTask_Run3(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	common.InitRunningEnv()
	logger.InitLogger()
	db.InitDB()
	db.DB.AutoMigrate(&Task{})
	db.DB.AutoMigrate(&TaskInstance{})
	db.DB.AutoMigrate(&TaskCmd{})
	db.DB.AutoMigrate(&TaskScript{})
	db.DB.AutoMigrate(&TaskStageAgent{})
	db.DB.AutoMigrate(&TaskArg{})
	db.DB.AutoMigrate(&InsArg{})

	ipList := make([]string, 0)
	for i := 0; i < 20; i++ {
		rands := rand_str.GetStr(4)
		ipList = append(ipList, rands)
	}

	collValByte, _ := json.Marshal(ipList)
	stage100 := []int{100}
	stg100, _ := json.Marshal(stage100)

	stage5050 := []int{50, 50}
	stg5050, _ := json.Marshal(stage5050)

	// *** ROOT TASK ***
	tr := Task{
		Name:            "test run3 task 28",
		TkType:          TASKROOT,
		CollectionType:  CollList,
		CollectionValue: string(collValByte),
		Stages:          string(stg100),
		Desc:            "a basic test task run3",
		//Stages:          string(sg), // 100%
		//ContinueOnFail: 1,
	}
	err = tr.Create()
	if err != nil {
		t.Fatalf("test create tr err: %s", err.Error())
	}

	// *** TEST CMD ***
	// *** stage: 100% ***
	// *** split: 50% ***
	t1 := Task{
		Name:            tr.Name,
		TkType:          TASKCMD,
		ParentTaskID:    tr.ID,
		SeqNum:          1,
		Stages:          string(stg100),
		ContinueOnFail:  1,
		SplitParent:     DoSplit,
		CollectionValue: strconv.Itoa(50),
	}
	err = t1.Create()
	if err != nil {
		t.Fatalf("test create t1 err: %s", err.Error())
	}
	c1 := TaskCmd{
		TaskID:  t1.ID,
		Cmd:     "echo1",
		Opt:     "123 ok",
		Timeout: 60,
	}
	err = c1.Create()
	if err != nil {
		t.Fatalf("test create c1 err: %s", err.Error())
	}

	// *** TEST SCRIPT ***
	// *** stage: 50% 50% ***
	// *** split: 50% ***
	t2 := Task{
		Name:            tr.Name,
		TkType:          TASKSCRIPT,
		ParentTaskID:    tr.ID,
		SeqNum:          2,
		Stages:          string(stg5050),
		SplitParent:     DoSplit,
		CollectionValue: strconv.Itoa(50),
	}
	err = t2.Create()
	if err != nil {
		t.Fatalf("test create t2 err: %s", err.Error())
	}
	s2 := TaskScript{
		TaskID:  t2.ID,
		Shell:   "/bin/csh2",
		Name:    "/tmp/test123.sh",
		RunAs:   "jacenr",
		Timeout: 60,
	}
	err = s2.Create()
	if err != nil {
		t.Fatalf("test create s2 err: %s", err.Error())
	}

	// *** TEST HGROUP ***
	// *** stage 40% 60% ***
	stageInt6040 := []int{40, 60}
	stg6040, _ := json.Marshal(stageInt6040)

	t3 := Task{
		Name:           tr.Name,
		TkType:         HGROUP,
		ParentTaskID:   tr.ID,
		SeqNum:         3,
		Stages:         string(stg6040), // 40% 60%
		ContinueOnFail: 1,
		//CollectionValue: strconv.Itoa(50),
		//SplitParent:     DoSplit,
	}
	err = t3.Create()
	if err != nil {
		t.Fatalf("test create t3 err: %s", err.Error())
	}

	// *** split 50% ***
	// *** test variables ***
	t31 := Task{
		Name:            tr.Name,
		TkType:          TASKCMD,
		ParentTaskID:    t3.ID,
		SeqNum:          1,
		CollectionValue: strconv.Itoa(50),
		Stages:          string(stg100), // 100%
		SplitParent:     DoSplit,
		ContinueOnFail:  1,
		TaskArgs: []TaskArg{
			{
				TaskName: tr.Name,
				ArgName:  "t31Var1",
				ArgType:  COMMONSTR,
				ArgValue: "t31Var1Value",
			},
			{
				TaskName: tr.Name,
				ArgName:  "t31Var2",
				ArgType:  AGENTFACT,
				ArgValue: "t31Var2Value",
			},
			{
				TaskName: tr.Name,
				ArgName:  "t31Var3",
				ArgType:  AGENTTAG,
				ArgValue: "t31Var3Value",
			},
		},
	}
	err = t31.Create()
	if err != nil {
		t.Fatalf("test create t31 err: %s", err.Error())
	}
	c31 := TaskCmd{
		TaskID:  t31.ID,
		Cmd:     "echo31",
		Opt:     "123 ok test",
		Timeout: 60,
	}
	err = c31.Create()
	if err != nil {
		t.Fatalf("test create c31 err: %s", err.Error())
	}

	t32 := Task{
		Name:            tr.Name,
		TkType:          TASKSCRIPT,
		ParentTaskID:    t3.ID,
		SeqNum:          2,
		CollectionValue: strconv.Itoa(50),
		Stages:          string(stg100), // 100%
		SplitParent:     DoSplit,
	}
	err = t32.Create()
	if err != nil {
		t.Fatalf("test create t32 err: %s", err.Error())
	}
	s32 := TaskScript{
		TaskID:  t32.ID,
		Shell:   "/bin/csh32",
		Name:    "/tmp/test123.sh",
		RunAs:   "jacenr",
		Timeout: 60,
	}
	err = s32.Create()
	if err != nil {
		t.Fatalf("test create s32 err: %s", err.Error())
	}

	// *** test continue depend task ***
	t33 := Task{
		Name:           tr.Name,
		TkType:         TASKSCRIPT,
		ParentTaskID:   t3.ID,
		SeqNum:         3,
		Stages:         string(stg100), // 100%
		SplitParent:    DoNotSplit,
		ContinueByTask: t31.ID,
		ContinueRst:    STAGEFAILED,
		ContinueOnFail: 1,
	}
	err = t33.Create()
	if err != nil {
		t.Fatalf("test create t33 err: %s", err.Error())
	}
	s33 := TaskScript{
		TaskID:  t33.ID,
		Shell:   "/bin/bash33",
		Name:    "/tmp/test.sh",
		RunAs:   "jacenr",
		Timeout: 60,
	}
	err = s33.Create()
	if err != nil {
		t.Fatalf("test create s33 err: %s", err.Error())
	}

	// *** test XGROUP ***
	t4 := Task{
		Name:           tr.Name,
		TkType:         XGROUP,
		ParentTaskID:   tr.ID,
		SeqNum:         4,
		Stages:         string(stg100),
		ContinueOnFail: 1,
		//CollectionValue: strconv.Itoa(50),
		//SplitParent:     DoSplit,
	}
	err = t4.Create()
	if err != nil {
		t.Fatalf("test create t4 err: %s", err.Error())
	}

	t41 := Task{
		Name:           tr.Name,
		TkType:         TASKCMD,
		ParentTaskID:   t4.ID,
		SeqNum:         1,
		Stages:         string(stg100),
		ContinueOnFail: 1,
	}
	err = t41.Create()
	if err != nil {
		t.Fatalf("test create t41 err: %s", err.Error())
	}
	c41 := TaskCmd{
		TaskID:  t41.ID,
		Cmd:     "echo41",
		Opt:     "123 ok",
		Timeout: 60,
	}
	err = c41.Create()
	if err != nil {
		t.Fatalf("test create c41 err: %s", err.Error())
	}

	t42 := Task{
		Name:           tr.Name,
		TkType:         TASKCMD,
		ParentTaskID:   t4.ID,
		SeqNum:         2,
		Stages:         string(stg100),
		ContinueOnFail: 1,
	}
	err = t42.Create()
	if err != nil {
		t.Fatalf("test create t42 err: %s", err.Error())
	}
	c42 := TaskCmd{
		TaskID:  t42.ID,
		Cmd:     "echo42",
		Opt:     "123 ok",
		Timeout: 60,
	}
	err = c42.Create()
	if err != nil {
		t.Fatalf("test create c42 err: %s", err.Error())
	}

	// *** test VGROUP ***
	// *** stage 100% ***
	t5 := Task{
		Name:           tr.Name,
		TkType:         VGROUP,
		ParentTaskID:   tr.ID,
		SeqNum:         5,
		Stages:         string(stg100),
		ContinueOnFail: 1,
	}
	err = t5.Create()
	if err != nil {
		t.Fatalf("test create t5 err: %s", err.Error())
	}

	t51 := Task{
		Name:           tr.Name,
		TkType:         TASKCMD,
		ParentTaskID:   t5.ID,
		SeqNum:         1,
		Stages:         string(stg100),
		ContinueOnFail: 1,
	}
	err = t51.Create()
	if err != nil {
		t.Fatalf("test create t51 err: %s", err.Error())
	}
	c51 := TaskCmd{
		TaskID:  t51.ID,
		Cmd:     "echo51",
		Opt:     "123 ok",
		Timeout: 60,
	}
	err = c51.Create()
	if err != nil {
		t.Fatalf("test create c51 err: %s", err.Error())
	}

	t52 := Task{
		Name:           tr.Name,
		TkType:         TASKCMD,
		ParentTaskID:   t5.ID,
		SeqNum:         2,
		Stages:         string(stg100),
		ContinueOnFail: 1,
	}
	err = t52.Create()
	if err != nil {
		t.Fatalf("test create t52 err: %s", err.Error())
	}
	c52 := TaskCmd{
		TaskID:  t52.ID,
		Cmd:     "echo52",
		Opt:     "123 ok",
		Timeout: 60,
	}
	err = c52.Create()
	if err != nil {
		t.Fatalf("test create c52 err: %s", err.Error())
	}

	t.Log("All tasks are created.")

	insName := NewInsName(tr.ID)
	ins := TaskInstance{
		Name: insName,
	}

	argMap := make(map[string]TaskArg)
	err = tr.InitInsArgs(argMap, insName)
	if err != nil {
		t.Fatalf("fail: %s", err)
	}

	err = tr.Run(&ins)
	if err != nil {
		t.Fatalf("test run tr fail: %s", err.Error())
	}

	allIns, err := GetAllInsByName(tr.Name)
	if err != nil {
		t.Fatalf("test run task fail: %s", err)
	}
	for _, v := range allIns {
		t.Logf("id:%d, parent_id:%d, task_id:%d, stage:%d, stage_seq:%d, start_at:%v, end_at:%v, status: %d, msg:%s\n",
			v.ID, v.ParentInsID, v.TaskID, v.Stage, v.StageSeq, v.StartAt, v.EndAt, v.Status, v.InsMsg)
	}

	time.Sleep(10 * time.Second)
}
