package task

import (
	"opsHeart_server/conf"
	"opsHeart_server/db"
	"testing"
)

func TestTask_Create(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&Task{})

	tk := Task{
		Name:            "a test task 1",
		TkType:          TASKROOT,
		CollectionType:  CollList,
		CollectionValue: "[\"A\",\"B\"]",
		CreateBy:        "jacen",
		Desc:            "just a test task",
	}

	err = tk.Create()
	if err != nil {
		t.Fatalf("test fail, create db data err: %s", err.Error())
	}

	tkc1 := Task{
		Name:         tk.Name,
		ParentTaskID: tk.ID,
		TkType:       TASKCMD,
		CreateBy:     "jacen",
		Desc:         "a child task",
		SeqNum:       1,
	}
	err = tkc1.Create()
	if err != nil {
		t.Fatalf("test fail, create tkc1 err: %s", err.Error())
	}

	tkc2 := Task{
		Name:         tk.Name,
		ParentTaskID: tk.ID,
		TkType:       TASKCMD,
		CreateBy:     "jacen",
		Desc:         "a child task",
		SeqNum:       2,
	}
	err = tkc2.Create()
	if err != nil {
		t.Fatalf("test fail, create tkc2 err: %s", err.Error())
	}
}

func TestTask_GetTheSeqChild(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&Task{})

	tk := Task{}
	tk.ID = 1

	var seq uint = 1
	tkc1, err := tk.GetTheSeqChild(seq)
	if err != nil {
		t.Fatalf("test get the seq %d child err: %s", seq, err.Error())
	}
	t.Logf("test success %v", tkc1)
}

func TestTask_GetAllChildTask(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&Task{})

	tk := Task{}
	tk.ID = 1

	allChild, err := tk.GetAllChildTask()
	if err != nil {
		t.Fatalf("tast gat all child task failed: %s", err)
	}

	for _, v := range allChild {
		t.Logf("test success, task: %v", v)
	}
}

func TestTask_GetBrotherTask(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&Task{})

	tk := Task{}
	tk.ParentTaskID = 1
	tk.SeqNum = 2

	tb, err := tk.GetOldBrotherTask()
	if err != nil {
		t.Fatalf("test get brother task failed: %s", err.Error())
	}
	t.Logf("test success, task id: %d, task seq: %d", tb.ID, tb.SeqNum)
}

func TestGetTaskByID(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&Task{})

	tk, err := GetTaskByID(1)
	if err != nil {
		t.Fatalf("test err: %s", err.Error())
	}
	t.Logf("test success, tk id: %d, tk name: %s", tk.ID, tk.Name)
}

func TestCheckNameIsExist(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&Task{})

	n1 := "a test name"
	n2 := "a test task 1"

	c1 := CheckNameIsExist(n1)
	t.Logf("test n1 rst: %v", c1)

	c2 := CheckNameIsExist(n2)
	t.Logf("test n1 rst: %v", c2)
}

func TestTask_GetNextTask(t *testing.T) {
	err := conf.InitCfg()
	if err != nil {
		t.Fatalf("init conf err: %s", err.Error())
	}
	db.InitDB()
	db.DB.AutoMigrate(&Task{})

	tk := Task{}
	tk.ParentTaskID = 1
	tk.SeqNum = 2

	next, err := tk.GetNextTask()
	if err != nil {
		t.Fatalf("test fail: %s", err.Error())
	}

	t.Logf("test success len: %d", len(next))
	for _, v := range next {
		t.Logf("test success: %v", v)
	}
}
