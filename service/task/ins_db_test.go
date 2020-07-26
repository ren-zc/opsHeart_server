package task

import (
	"opsHeart_server/conf"
	"opsHeart_server/db"
	"testing"
	"time"
)

func TestInsArg_Create(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&InsArg{})

	ia := InsArg{
		InsName: "a test ins name 999",
		TaskID:  99999999,
		ArgName: "TestArgName999",
		ArgType: COMMONSTR,
	}

	err = ia.Create()
	if err != nil {
		t.Fatalf("fail: %s", err)
	}
	t.Log("success!")
}

func TestTaskInstance_Create(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&TaskInstance{})

	ins1 := TaskInstance{
		Name:        NewInsName(12),
		TaskID:      12,
		StageSeq:    1,
		StageAgents: "12-100-vUNq",
		StartAt:     time.Now(),
	}
	err = ins1.Create()
	if err != nil {
		t.Fatalf("test create instance fail: %s", err.Error())
	}
	t.Logf("success: %v", ins1)

	ins2 := TaskInstance{
		Name:        ins1.Name,
		TaskID:      2,
		StageSeq:    1,
		StageAgents: "12-100-vUNq",
		StartAt:     time.Now(),
	}
	ins2.ParentInsID = ins1.ID
	err = ins2.Create()
	if err != nil {
		t.Fatalf("test create instance fail: %s", err.Error())
	}
	t.Logf("success: %v", ins2)

	ins3 := TaskInstance{
		Name:        ins1.Name,
		TaskID:      3,
		StageSeq:    1,
		StageAgents: "12-100-vUNq",
		StartAt:     time.Now(),
	}
	ins3.ParentInsID = ins1.ID
	err = ins3.Create()
	if err != nil {
		t.Fatalf("test create instance fail: %s", err.Error())
	}
	t.Logf("success: %v", ins3)
}

func TestTaskInstance_Update(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&TaskInstance{})

	ins := TaskInstance{}
	ins.ID = 1

	err = ins.Update([]string{"status", "start_at"}, STAGERUNNING, time.Now())
	if err != nil {
		t.Fatalf("test update fail: %s", err.Error())
	}

	t.Logf("success: %v", ins)
}

func TestGetAllInsByNameAndTaskID(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&TaskInstance{})

	all, err := GetAllInsByNameAndTaskID("12-20206230516-efRX", 2)
	if err != nil {
		t.Fatalf("get all instance fail: %s", err.Error())
	}
	t.Logf("test success: %v", all)
}

func TestTaskInstance_GetNextInstance(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&TaskInstance{})

	ins := TaskInstance{
		ParentInsID: 1,
		TaskID:      3,
		StageSeq:    0,
	}

	next, err := ins.GetNextStageInstance()
	if err != nil {
		t.Fatalf("test fail: %s", err.Error())
	}

	//t.Logf("test success: %v", next)
	t.Logf("test next len: %d", len(next))
	for _, v := range next {
		t.Logf("test success: %v", v)
	}
}
