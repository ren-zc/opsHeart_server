package server

import (
	"fmt"
	"opsHeart_server/common"
	"opsHeart_server/conf"
	"opsHeart_server/logger"
	"opsHeart_server/routers"
)

type Data struct {
	Version     string
	BuildTime   string
	GoVersion   string
	Author      string
	ReleaseTime string
}

var Svr Data

func init() {
	Svr.Version = common.Version
	Svr.BuildTime = common.BuildTime
	Svr.Author = common.Author
	Svr.ReleaseTime = common.ReleaseTime
}

func Start() {
	ip := conf.GetAddr()
	port := conf.GetPort()
	addr := fmt.Sprintf("%s:%s", ip, port)
	err := routers.R.Run(addr)
	if err != nil {
		logger.ServerLog.Errorf("action=run http server;err=%s", err.Error())
	}
}
