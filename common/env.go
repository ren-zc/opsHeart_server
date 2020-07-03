package common

import "opsHeart/conf"

var AgentPort string

func InitAgentPort() {
	AgentPort = conf.GetAgentPort()
}
