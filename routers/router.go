package routers

import (
	"github.com/gin-gonic/gin"
	"opsHeart_server/handlers/agent/start_up"
	agentV1 "opsHeart_server/handlers/agent/v1"
	externalV1 "opsHeart_server/handlers/external/v1"
	frontV1 "opsHeart_server/handlers/front/v1"
	testV1 "opsHeart_server/handlers/test/v1"
	"opsHeart_server/routers/middleware"
)

var R *gin.Engine

func init() {
	R = gin.Default()

	// for agent call
	agent := R.Group("/agent")
	{
		v1 := agent.Group("/v1")
		v1.Use(middleware.TokenChecker())
		{
			v1.GET("/hbs", agentV1.HandleHbs)
			v1.GET("/status", agentV1.HandleStatus)
			v1.POST("/fact", agentV1.HandleFact)
		}
		startUp := agent.Group("/start-up")
		{
			startUp.POST("/register", start_up.HandleAgentRegister)
		}
	}

	// for front call
	front := R.Group("/front")
	front.Use(middleware.FrontToken())
	{
		v1 := front.Group("/v1")
		{
			v1.GET("/register-agents", frontV1.HandleQueryUnregAgents)
			v1.POST("/register-agents", frontV1.HandleStartAgents)

			// manage token of external
			v1.POST("/manage-external", nil)

			// run task
			v1.POST("/task/run", frontV1.HandleRunTask)
		}
		//col := v1.Group("/collection")
	}

	// for external system call
	external := R.Group("/external")
	{
		v1 := external.Group("/v1")
		{
			v1.GET("/test", nil)
			v1.POST("/task/run", externalV1.HandleRunTask)
		}
	}

	// for test
	test := R.Group("/test")
	{
		v1 := test.Group("/v1")
		{
			v1.GET("/do", testV1.TestGetHandler)
			v1.POST("/do", testV1.TestPostHandler)
		}
	}

	//external := R.Group("/external")
	//{
	//	v1 := external.Group("/v1")
	//}
}
