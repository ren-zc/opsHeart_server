package cron

import "reflect"

// IsFunc check v if a func
func IsFunc(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Func
}

// [ NOT USED ] Invoke call func interface
func Invoke(fn interface{}, args ...interface{}) {
	v := reflect.ValueOf(fn)
	rargs := make([]reflect.Value, len(args))
	for i, a := range args {
		rargs[i] = reflect.ValueOf(a)
	}
	v.Call(rargs)
}
