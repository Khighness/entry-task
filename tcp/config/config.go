package config

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/Khighness/entry-task/tcp/logging"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-04-28

type AppConfig struct {
	Server *ServerConfig `yaml:"server"`
	MySQL  *MySQLConfig  `yaml:"mysql"`
	Redis  *RedisConfig  `yaml:"redis"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type MySQLConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Pass        string `yaml:"pass"`
	Name        string `yaml:"name"`
	MaxOpen     int    `yaml:"max-open"`
	MaxIdle     int    `yaml:"max-idle"`
	MaxLifeTime int    `yaml:"max-life-time"`
	MaxIdleTime int    `yaml:"max-idle-time"`
}

type RedisConfig struct {
	Addr        string `yaml:"addr"`
	Pass        string `yaml:"pass"`
	Db          int    `yaml:"db"`
	MaxConn     int    `yaml:"max-conn"`
	MinIdle     int    `yaml:"min-idle"`
	DialTimeout int    `yaml:"dial-timeout"`
	IdleTimeout int    `yaml:"idle-timeout"`
	MaxRetries  int    `yaml:"max-retries"`
	MaxConnAge  int    `yaml:"max-conn-age"`
}

var AppCfg *AppConfig

var (
	port = flag.Int("p", 20000, "Web service port")
)

// Load 导入配置和参数
func Load() {
	AppCfg = &AppConfig{}
	applicationFile, err := ioutil.ReadFile("application-tcp.yml")
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
}
