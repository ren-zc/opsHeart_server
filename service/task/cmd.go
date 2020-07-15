package task

import (
	"fmt"
	"opsHeart_server/logger"
)

func (c *TaskCmd) start(agents *[]string) error {
	//return errors.New("test")
	logger.TaskLog.Debugf("CMD run, cmd:%s, opt:%s.\n", c.Cmd, c.Opt)
	for _, v := range *agents {
		fmt.Printf("\t%s\n", v)
	}
	return nil
}
