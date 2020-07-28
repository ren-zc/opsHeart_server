package common

import "opsHeart_server/conf"

var AgentPort string
var StepPause bool

func InitRunningEnv() {
	AgentPort = conf.GetAgentPort()
	StepPause = conf.GetStepPause()
}
