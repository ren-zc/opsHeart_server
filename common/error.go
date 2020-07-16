package common

type ErrorCode uint

const (
	// 0-999
	NoError ErrorCode = 0

	// 1000-1999
	BindPostDataErr ErrorCode = 1000

	// 2000-2999
	RunTaskErr           ErrorCode = 2000
	QueryTaskByIDErr     ErrorCode = 2001
	QueryInstanceByIDErr ErrorCode = 2002
)

var ErrCodeDict = map[ErrorCode]string{
	BindPostDataErr:      "bind post data error",
	RunTaskErr:           "run task err",
	QueryTaskByIDErr:     "query task by id error",
	QueryInstanceByIDErr: "query instance by id error",
}
