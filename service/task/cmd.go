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
		fmt.Printf("args: %v\n", c.Args)
		for _, a := range c.Args {
			fmt.Printf("\targ name: %s, arg type: %d, arg value: %s\n",
				a.ArgName, a.ArgType, a.ArgValue)
		}
	}
	return nil
}
