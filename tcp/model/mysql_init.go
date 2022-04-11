package model

import (
	"database/sql"
	"github.com/Khighness/entry-task/tcp/logging"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
	"strings"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

var (
	DB     *sql.DB
	DbHost string
	DbPort string
	DbUser string
	DbPass string
	DbName string
)

// Load 初始化MySQL
func Load(file *ini.File) {
	loadMysqlConfig(file)
	url := strings.Join([]string{DbUser, ":", DbPass, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
	connectMySQL(url)
}

// loadMysqlConfig 读取MySQL配置
func loadMysqlConfig(file *ini.File) {
	mysql := file.Section("mysql")
	DbHost = mysql.Key("DbHost").String()
	DbPort = mysql.Key("DbPort").String()
	DbUser = mysql.Key("DbUser").String()
	DbPass = mysql.Key("DbPass").String()
	DbName = mysql.Key("DbName").String()
}

// connectMySQL 连接到MySQL
func connectMySQL(url string) {
	db, err := sql.Open("mysql", url)

	if err != nil {
		logging.Log.Fatalf("Wrong configuration of [MySQL] in config file: %s", err)
	}

	// 最大连接数
	db.SetMaxOpenConns(10000)
	// 闲置连接数
	db.SetMaxIdleConns(1000)
	// 最大存活时间
	db.SetConnMaxLifetime(time.Hour)
	// 最大空闲时间
	db.SetConnMaxIdleTime(time.Hour)

	DB = db
	if err = DB.Ping(); err != nil {
		logging.Log.Fatalf("Failed to connect to mysql server [%s]: %s", strings.Join([]string{DbHost, ":", DbPort}, ""), err)
	} else {
		logging.Log.Infof("Succeed to connect to mysql server [%s]", strings.Join([]string{DbHost, ":", DbPort}, ""))
	}
}
