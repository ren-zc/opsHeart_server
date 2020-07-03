package worker

import (
	"fmt"
	"testing"
	"time"
)

func TestNewWorker(t *testing.T) {
	printStr := func(s string) {
		time.Sleep(1 * time.Second)
		fmt.Println(s)
	}

	w1 := NewWorker(10)
	w1.StartWork()
	defer w1.EndWorkerAndWait()
	//defer w1.ForceEndWorker()
	for i := 'A'; i <= 'z'; i++ {
		func(n string) {
			w1.Ch <- func() {
				printStr(n)
			}
		}(string(i))
	}
}
