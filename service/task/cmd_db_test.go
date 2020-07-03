package task

import (
	"opsHeart/conf"
	"opsHeart/db"
	"testing"
)

func TestTaskCmd_QueryByTaskID(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&TaskCmd{})

	c := TaskCmd{
		TaskID: 2,
		Cmd:    "echo test",
	}
	err = c.QueryByTaskID()
	if err != nil {
		t.Fatalf("test query cmd task by task id err: %s", err.Error())
	}

	t.Logf("task id:%d, cmd:%s, opt:%s, timeout:%d",
		c.TaskID, c.Cmd, c.Opt, c.Timeout)
}
