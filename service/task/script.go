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
	logger.TaskLog.Debugf("script run, shell:%s, name:%s.\n", s.Shell, s.Name)
	for _, v := range *agents {
		fmt.Printf("\t%s\n", v)
	}
	return nil
}
