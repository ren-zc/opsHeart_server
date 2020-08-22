package task

import (
	"errors"
	"fmt"
	"opsHeart_server/logger"
	"strings"
)

func (s *TaskScript) start(agents *[]string) error {
	if strings.Contains(s.Shell, "zsh") {
		return errors.New("test zsh")
	}
	logger.TaskLog.Infof("script run, shell:%s, name:%s.\n", s.Shell, s.Name)
	for _, v := range *agents {
		fmt.Printf("\t%s\n", v)
		for _, a := range s.Args {
			fmt.Printf("\targ name: %s, arg type: %d, arg value: %s\n",
				a.ArgName, a.ArgType, a.ArgValue)
		}
	}
	return nil
}
