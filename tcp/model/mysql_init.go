package model

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Khighness/entry-task/tcp/config"
	"github.com/Khighness/entry-task/tcp/logging"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

var (
	DB *sql.DB
)

// InitMySQL 初始化MySQL连接池
func InitMySQL() {
	DB = ConnectMySQL(config.AppCfg.MySQL)
}

// ConnectMySQL 连接到MySQL
func ConnectMySQL(mysqlCfg *config.MySQLConfig) *sql.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", mysqlCfg.User, mysqlCfg.Pass, mysqlCfg.Host, mysqlCfg.Port, mysqlCfg.Name, "charset=utf8&parseTime=true")
	db, err := sql.Open("mysql", url)
	if err != nil {
		logging.Log.Fatalf("Wrong configuration of [MySQL] in config file: %s", err)
	}

	// 最大连接数
	db.SetMaxOpenConns(mysqlCfg.MaxOpen)
	// 闲置连接数
	db.SetMaxIdleConns(mysqlCfg.MaxIdle)
	// 最大存活时间
	db.SetConnMaxLifetime(time.Duration(mysqlCfg.MaxLifeTime) * time.Second)
	// 最大空闲时间
	db.SetConnMaxIdleTime(time.Duration(mysqlCfg.MaxIdleTime) * time.Second)

	if err = db.Ping(); err != nil {
		logging.Log.Fatalf("Failed to connect to mysql server [%s]: %s", fmt.Sprintf("%s:%d", mysqlCfg.Host, mysqlCfg.Port), err)
	} else {
		logging.Log.Infof("Succeed to connect to mysql server [%s]", fmt.Sprintf("%s:%d", mysqlCfg.Host, mysqlCfg.Port))
	}
	return db
}
