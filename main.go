package main

import (
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"opsHeart_server/common"
	"opsHeart_server/conf"
	"opsHeart_server/db"
	"opsHeart_server/logger"
	"opsHeart_server/server"
	"opsHeart_server/service/agent"
	"opsHeart_server/service/collection"
	"opsHeart_server/service/cron_task"
	"os"
	"time"
)

var (
	h         bool
	v         bool
	u         bool
	buildTime string
	goVersion string
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.BoolVar(&v, "v", false, "show version, build info.")
	flag.BoolVar(&u, "u", false, "create uuid")
	flag.Usage = usage
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `Options: Build Time: %v`, buildTime)
	flag.PrintDefaults()
}

func main() {
	common.BuildTime = buildTime
	common.GoVersion = goVersion

	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	if v {
		fmt.Println("Proxy Version: ", common.Version)
		fmt.Println("Build Time: ", common.BuildTime)
		fmt.Println("Go Version: ", common.GoVersion)
		return
	}

	if u {
		uuidObj := uuid.NewV4()
		uu, err := uuidObj.Value()

		if err != nil {
			fmt.Printf("create uuid err: %s\n", err)
		} else {
			fmt.Printf("uuid is: %s\n", uu)
		}
		return
	}

	// init conf
	err := conf.InitCfg()
	if err != nil {
		fmt.Printf("action=init conf;err=%s\n", err)
		return
	}

	logger.ServerLog.Info("action=load conf file;status=success")

	common.InitAgentPort()

	logger.InitLogger()

	db.InitDB()
	defer func() {
		_ = db.DB.Close()
	}()

	// init model
	db.DB.AutoMigrate(agent.Agent{})
	db.DB.AutoMigrate(collection.AgentFact{})

	if conf.GetServerRole() {
		// start hbs check cron task
		go func() {
			time.Sleep(3 * time.Minute)
			err := cron_task.CheckHbs.Start()
			if err != nil {
				logger.HbsLog.Errorf("action=start check hbs cron;err=%s", err.Error())
			}
		}()
		defer func() {
			_ = cron_task.CheckHbs.Stop()
		}()
	}

	server.Start()
}
