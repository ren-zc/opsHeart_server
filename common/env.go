package common

import "opsHeart_server/conf"

var AgentPort string

func InitAgentPort() {
	AgentPort = conf.GetAgentPort()
}
