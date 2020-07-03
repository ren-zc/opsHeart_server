package conf

import "github.com/go-ini/ini"

var cfg *ini.File

func InitCfg() error {
	var err error
	cfg, err = ini.Load("./conf/app.ini")
	return err
}

// GetPort get port config
func GetPort() string {
	return cfg.Section("server").Key("http_port").String()
}

// GetAddr get addr config
func GetAddr() string {
	return cfg.Section("server").Key("http_addr").String()
}

func GetLogLevel() string {
	return cfg.Section("app").Key("log_level").String()
}

func GetMode() string {
	return cfg.Section("app").Key("mode").String()
}

func GetDbType() string {
	return cfg.Section("db").Key("type").String()
}

func GetDbUrl() string {
	return cfg.Section("db").Key("url").String()
}

func GetUUID() string {
	return cfg.Section("server").Key("uuid").String()
}

func GetFrontToken() string {
	return cfg.Section("server").Key("front_token").String()
}

func GetAgentPort() string {
	return cfg.Section("server").Key("agent_port").String()
}

func GetServerRole() bool {
	return cfg.Section("cluster").Key("role").MustInt() == 1
}

func GetWorkerConcurrent() int {
	return cfg.Section("server").Key("worker_concurrent").MustInt()
}
