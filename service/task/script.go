package task

import (
	"errors"
	"fmt"
	"strings"
)

func (s *TaskScript) start(agents *[]string) error {
	if strings.Contains(s.Shell, "zsh") {
		return errors.New("test zsh")
	}
	fmt.Printf("script run, shell:%s, name:%s.\n", s.Shell, s.Name)
	for _, v := range *agents {
		fmt.Printf("\t%s\n", v)
	}
	return nil
}
