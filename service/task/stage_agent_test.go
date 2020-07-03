package task

import (
	"opsHeart_server/conf"
	"opsHeart_server/db"
	"opsHeart_server/utils/rand_str"
	"testing"
)

func TestTaskStageAgent_Create(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&TaskStageAgent{})
	n := 10
	sn := NewStageName(12, 100)
	for i := 0; i < n; i++ {
		sa := TaskStageAgent{
			StageName: sn,
			IP:        rand_str.GetStr(4),
		}
		err := sa.Create()
		if err != nil {
			t.Fatalf("create sa: %v, err: %s, n: %d", sa, err.Error(), n)
		}
	}
}

func TestTaskStageAgent_Update(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&TaskStageAgent{})

	var allSa []TaskStageAgent
	err = db.DB.Model(&TaskStageAgent{}).Where("child_use = ?", IsUsed).Find(&allSa).Error
	if err != nil {
		t.Fatalf("query all sa err: %s", err.Error())
	}

	for _, v := range allSa {
		err := v.Update("child_use", NotUsed)
		if err != nil {
			t.Fatalf("update sa err: %s", err)
		}
	}
}
