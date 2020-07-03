package cron

import (
	"context"
	"errors"
	"reflect"
	"time"
)

type Status int

const (
	STARTED Status = 1
	STOPPED Status = 0
)

type Cr struct {
	f    interface{}   // the func called by period
	args []interface{} // func arguments
	d    time.Duration // the duration
	c    int           // the max concurrent
	r    int           // number in running
	s    Status        // cron status
	//ch     chan int  // handle r
	ctx    context.Context
	cancel context.CancelFunc
}

// init cron object
func NewCron(f interface{}, args []interface{}, d time.Duration, c int) (*Cr, error) {
	if !IsFunc(f) {
		return nil, errors.New("func is not callable function")
	}

	o := Cr{
		f:    f,
		args: args,
		d:    d,
		c:    c,
	}

	ctxB := context.Background()
	ctx, cancel := context.WithCancel(ctxB)

	o.ctx = ctx
	o.cancel = cancel

	return &o, nil
}

// start cron object
func (c *Cr) Start() error {
	if c.s == STARTED {
		return errors.New("cron had started, do not start again")
	}
	c.s = STARTED

	// init a channel for manager of Cr.r
	//ch := make(chan int)
	//c.ch = ch
	//go c.handleConcurrentNum()

	go func() {
		t := time.NewTicker(c.d)
		defer t.Stop()

		f := reflect.ValueOf(c.f)
		rargs := make([]reflect.Value, len(c.args))
		for i, a := range c.args {
			rargs[i] = reflect.ValueOf(a)
		}

		fc := func() {
			//fmt.Println(c.r)  // test
			// may need to be optimized:
			// locker, token channel or a goroutine whose role is a manager
			c.r++
			f.Call(rargs)
			c.r--
		}

		go fc()

		for {
			select {
			case <-c.ctx.Done():
				return
			case <-t.C:
				if c.r >= c.c {
					continue
				}
				go fc()
			}
		}
	}()
	return nil
}

// stop the cron
func (c *Cr) Stop() error {
	if c.s == STOPPED {
		return errors.New("cron had stop, do not stop again")
	}
	c.cancel()
	c.s = STOPPED
	return nil
}

// A goroutine whose role is manager of Cr.r
//func (c *Cr) handleConcurrentNum() {
//	for {
//		i, ok := <-c.ch
//		if !ok {
//			break
//		}
//		c.r += i
//	}
//}
