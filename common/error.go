package common

type ErrorCode uint

const (
	// 1000-1999
	BindPostDataErr ErrorCode = 1000
	// 2000-2999
	RunTaskErr ErrorCode = 2000
)

var ErrCodeDict = map[ErrorCode]string{
	BindPostDataErr: "bind post data error",
	RunTaskErr:      "run task err",
}
