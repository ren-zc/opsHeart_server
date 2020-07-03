package task

import (
	"encoding/json"
	"opsHeart/conf"
	"opsHeart/db"
	"opsHeart/utils/rand_str"
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
	db.InitDB()
	db.DB.AutoMigrate(&Task{})
	db.DB.AutoMigrate(&TaskInstance{})
	db.DB.AutoMigrate(&TaskCmd{})
	db.DB.AutoMigrate(&TaskScript{})
	db.DB.AutoMigrate(&TaskStageAgent{})

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
		Name:            "test task 50",
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
	}
	err = t1.Create()
	if err != nil {
		t.Fatalf("test create t1 err: %s", err.Error())
	}
	//s1 := TaskScript{
	//	TaskID:  t1.ID,
	//	Shell:   "/bin/t1sh",
	//	Name:    "/tmp/test123.sh",
	//	RunAs:   "jacenr",
	//	Timeout: 60,
	//}
	//err = s1.Create()
	//if err != nil {
	//	t.Fatalf("test create s1 err: %s", err.Error())
	//}
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

	//t21 := Task{
	//	Name:            tr.Name,
	//	TkType:          TASKCMD,
	//	ParentTaskID:    t2.ID,
	//	SeqNum:          1,
	//	CollectionValue: strconv.Itoa(50),
	//	Stages:          string(sg),
	//	SplitParent:     DoSplit,
	//	ContinueOnFail:  1,
	//}
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
	//s21 := TaskScript{
	//	TaskID:  t21.ID,
	//	Shell:   "/bin/t21sh",
	//	Name:    "/tmp/test123.sh",
	//	RunAs:   "jacenr",
	//	Timeout: 60,
	//}
	//err = s21.Create()
	//if err != nil {
	//	t.Fatalf("test create s21 err: %s", err.Error())
	//}
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
	}
	err = t31.Create()
	if err != nil {
		t.Fatalf("test create t31 err: %s", err.Error())
	}
	//s31 := TaskScript{
	//	TaskID:  t31.ID,
	//	Shell:   "/bin/zsh",
	//	Name:    "/tmp/test123.sh",
	//	RunAs:   "jacenr",
	//	Timeout: 60,
	//}
	//err = s31.Create()
	//if err != nil {
	//	t.Fatalf("test create s31 err: %s", err.Error())
	//}
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

	ins := TaskInstance{
		Name: NewInsName(tr.ID),
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
