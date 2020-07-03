package cron

import (
	"fmt"
	"testing"
	"time"
)

func printNS(i int64, s string) {
	time.Sleep(time.Duration(i) * time.Second)
	fmt.Printf("n is: %d, s is: %s\n", i, s)
}

func printFunc() {
	time.Sleep(time.Duration(1) * time.Second)
	fmt.Printf("this func no argments.\n")
}

// test cron use
func TestCr_Start1(t *testing.T) {
	argsList := []interface{}{int64(2), "a"}
	c, err := NewCron(printNS, argsList, 3*time.Second, 1)
	if err != nil {
		t.Fatalf("fail: %s\n", err.Error())
	}
	_ = c.Start()
	time.Sleep(15 * time.Second)
	_ = c.Stop()
	time.Sleep(5 * time.Second) // wait last task complete.
}

// test cron with func which no argments
func TestCr_Start2(t *testing.T) {
	c, err := NewCron(printFunc, nil, 3*time.Second, 1)
	if err != nil {
		t.Fatalf("fail: %s\n", err.Error())
	}
	_ = c.Start()
	time.Sleep(15 * time.Second)
	_ = c.Stop()
	time.Sleep(5 * time.Second) // wait last task complete.
}

// test start twice and stop twice
func TestCr_Stop(t *testing.T) {
	argsList := []interface{}{int64(3), "a"}
	c, err := NewCron(printNS, argsList, 2*time.Second, 1)
	if err != nil {
		t.Fatalf("fail: %s\n", err.Error())
	}

	_ = c.Start()
	err = c.Start()
	if err != nil {
		t.Logf("start err: %s\n", err.Error())
	}

	time.Sleep(15 * time.Second)

	_ = c.Stop()
	err = c.Stop()
	if err != nil {
		t.Logf("stop err: %s\n", err.Error())
	}

	time.Sleep(5 * time.Second)
}
