package agent

import (
	"opsHeart_server/conf"
	"opsHeart_server/db"
	"testing"
)

func TestGetAllUnreg(t *testing.T) {
	_ = conf.InitCfg()
	db.InitDB()
	all, err := GetAllUnreg()
	if err != nil {
		t.Fatalf("fail: %s", err.Error())
	}
	t.Logf("all: %v", all)
}
