package task

import (
	"fmt"
)

func (c *TaskCmd) start(agents *[]string) error {
	//return errors.New("test")
	fmt.Printf("CMD run, cmd:%s, opt:%s.\n", c.Cmd, c.Opt)
	for _, v := range *agents {
		fmt.Printf("\t%s\n", v)
	}
	return nil
}
