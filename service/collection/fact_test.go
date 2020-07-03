package collection

import (
	"opsHeart/conf"
	"opsHeart/db"
	"testing"
)

func initAf() error {
	err := conf.InitCfg()
	if err != nil {
		return err
	}
	db.InitDB()
	db.DB.AutoMigrate(AgentFact{})
	return nil
}

func TestAgentFact_Create(t *testing.T) {
	err := initAf()
	if err != nil {
		t.Fatalf("test create af, init cfg err: %s", err)
	}

	// The 2nd test will be failed.
	af1 := AgentFact{
		UUID:  "test-uuid1",
		Key:   "test-key2",
		Value: "test-value3",
	}

	err = af1.Create()
	if err != nil {
		t.Fatalf("test create af1 fail: %s", err.Error())
	}
	t.Logf("test create af1 success!")

	// The 2nd test will be failed.
	af2 := AgentFact{
		UUID:  "test-uuid2",
		Key:   "test-key2",
		Value: "test-value2",
	}

	err = af2.Create()
	if err != nil {
		t.Fatalf("test create af2 err: %s", err.Error())
	}
	t.Logf("test create af2 success!")
}

func TestAgentFact_IsExist(t *testing.T) {
	err := initAf()
	if err != nil {
		t.Fatalf("test create af, init cfg err: %s", err)
	}

	af1 := AgentFact{
		UUID:  "test-uuid1",
		Key:   "test-key1",
		Value: "test-value1",
	}
	af2 := AgentFact{
		UUID:  "test-uuid2",
		Key:   "test-key2",
		Value: "test-value2",
	}

	af1IsExist := af1.IsExist()
	if af1IsExist {
		t.Log("af1 is exist, af1 test ok!")
	} else {
		t.Error("af2 is not exist, af1 test fail!")
	}

	af2IsExist := af2.IsExist()
	if af2IsExist {
		t.Log("af2 is exist, af2 test ok!")
	} else {
		t.Error("af2 is not exist, af2 test fail!")
	}
}

func TestAgentFact_Update(t *testing.T) {
	err := initAf()
	if err != nil {
		t.Fatalf("test create af, init cfg err: %s", err)
	}

	af1 := AgentFact{
		UUID:  "test-uuid1",
		Key:   "test-key1",
		Value: "test-value2",
	}

	err = af1.Update()
	if err != nil {
		t.Fatalf("test update err: %s", err)
	}

	// test query by the way.
	s, err := af1.QueryValue()
	if err != nil {
		t.Fatalf("test update and query err: %s", err)
	}
	t.Logf("test update query value: %s", s)
}

func TestAgentFact_QueyAllKeyValueByUUID(t *testing.T) {
	err := initAf()
	if err != nil {
		t.Fatalf("test create af, init cfg err: %s", err)
	}

	af1 := AgentFact{
		UUID: "test-uuid1",
		Key:  "test-key1",
	}

	afs, err := af1.QueyAllKeyValueByUUID()

	t.Logf("afs length: %d", len(afs))

	if err != nil {
		t.Fatalf("test query all by uuid err: %s", err.Error())
	}
	for _, a := range afs {
		t.Logf("uuid:%s, key:%s, value: %s", a.UUID, a.Key, a.Value)
	}
}

func TestAgentFact_DeleteAKey(t *testing.T) {
	err := initAf()
	if err != nil {
		t.Fatalf("test create af, init cfg err: %s", err)
	}

	af2 := AgentFact{
		UUID:  "test-uuid2",
		Key:   "test-key2",
		Value: "test-value2",
	}

	err = af2.DeleteAKey()
	if err != nil {
		t.Fatalf("test delete a key err: %s", err.Error())
	}
	t.Log("test delete a key success!")
}

func TestAgentFact_DeleteUUIDAll(t *testing.T) {
	err := initAf()
	if err != nil {
		t.Fatalf("test create af, init cfg err: %s", err)
	}

	af1 := AgentFact{
		UUID:  "test-uuid1",
		Key:   "test-key1",
		Value: "test-value2",
	}

	err = af1.DeleteUUIDAll()
	if err != nil {
		t.Fatalf("test delete all by uuid err: %s", err)
	}
	t.Logf("test delete all by uuid success!")
}
