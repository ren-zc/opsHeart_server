package task

import (
	"opsHeart_server/conf"
	"opsHeart_server/db"
	"testing"
)

func TestTaskInstance_GetIPsByStageName(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	ins := TaskInstance{}
	db.DB.Model(&TaskInstance{}).Where("stage_agents = ?", "12-100-vUNq").
		First(&ins)
	t.Logf("ins id: %d, ins stage_agents: %s", ins.ID, ins.StageAgents)
	ips, err := ins.GetAllInsIPs()
	if err != nil {
		t.Fatalf("test get ips fail: %s", err.Error())
	}
	t.Logf("success: %v", ips)
}

func TestTaskInstance_GetAllInsIPsByPercent(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	ins := TaskInstance{}
	db.DB.Model(&TaskInstance{}).Where("stage_agents = ?", "12-100-vUNq").
		First(&ins)
	t.Logf("ins id: %d, ins stage_agents: %s", ins.ID, ins.StageAgents)
	//ips, err := ins.GetAllInsIPsByPercent(2, StageNumber)
	ips, err := ins.GetAllInsIPsByPercentOrNum(50, StagePercent, DoSplit)
	if err != nil {
		t.Fatalf("test get ips fail: %s", err.Error())
	}
	t.Logf("success: %v", ips)
}

func TestTaskInstance_GetAllStagesOfBrotherTask(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&TaskInstance{})

	ins := TaskInstance{}
	ins.TaskID = 3
	ins.ParentInsID = 1

	tk := Task{}
	tk.ParentTaskID = 1
	tk.SeqNum = 2

	allBt, err := ins.getAllStagesOfBrotherTask(&tk)
	if err != nil {
		t.Fatalf("test failed: %s", err.Error())
	}

	t.Logf("test success: %v", allBt)
}
