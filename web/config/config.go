package config

import (
	"github.com/Khighness/entry-task/web/logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
	MaxOpen     int    `yaml:"max-open"`
	MaxIdle     int    `yaml:"max-idle"`
	MaxLiftTime int    `yaml:"max-lift-time"`
	MaxIdleTime int    `yaml:"max-idle-time"`
}

var AppCfg *AppConfig

func init() {
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
}
