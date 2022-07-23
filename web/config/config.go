package config

import (
	"flag"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/Khighness/entry-task/web/logging"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-04-28

type AppConfig struct {
	Server *ServerConfig `yaml:"server"`
	Rpc    *RpcConfig    `yaml:"rpc"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type RpcConfig struct {
	Addr        string `yaml:"addr"`
	Initial     int    `yaml:"initial"`
	MaxOpen     int    `yaml:"max-open"`
	MaxIdle     int    `yaml:"max-idle"`
	MaxLifeTime int    `yaml:"max-life-time"`
	MaxIdleTime int    `yaml:"max-idle-time"`
}

var AppCfg *AppConfig

var (
	port    = flag.Int("p", 10000, "Web service port")
	rpcPort = flag.Int("rp", 20000, "Rpc service port")
)

// Load 导入配置和参数
func Load() {
	AppCfg = &AppConfig{}
	applicationFile, err := ioutil.ReadFile("application-web.yml")
	if err != nil {
		logging.Log.Fatal("Failed to load application configuration file")
	}
	err = yaml.Unmarshal(applicationFile, AppCfg)
	if err != nil {
		logging.Log.Fatal("Failed to read application configuration file")
	}
	logging.Log.Println("Succeed to load application configuration file")

	flag.Parse()
	AppCfg.Server.Port = *port
	AppCfg.Rpc.Addr = fmt.Sprintf("0.0.0.0:%v", *rpcPort)
}
