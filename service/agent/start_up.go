package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"opsHeart_server/common"
	"opsHeart_server/utils/call_http"
)

type tokenMsg struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Token  string `json:"token"`
}

func (a *Agent) StartUpAgent() error {
	host := fmt.Sprintf("%s:%s", a.RemoteAddr, common.AgentPort)
	var d tokenMsg
	if a.Status == ACCEPTED {
		d.Msg = "accepted"
		d.Status = 1
		d.Token = a.Token
	} else {
		d.Msg = "denied"
		d.Status = -1
	}

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	c, _, err := call_http.HttpPost(host, a.Token, common.AgentStartUpPath, b)
	if err != nil {
		return err
	}

	if c != 200 {
		return errors.New("agent resp code is not 200")
	}
	return nil
}
