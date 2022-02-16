package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
	"log"
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

// for test
func init()  {
	url := "root:KAG1823@tcp(127.0.0.1:3306)/entry_task?charset=utf8&parseTime=true"
	connectMySQL(url)
}

// MySQL 初始化MySQL
func MySQL(file *ini.File) {
	loadMysqlConfig(file)
	url := strings.Join([]string{DbUser, ":", DbPass, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
	connectMySQL(url)
}

// 读取MySQL配置
func loadMysqlConfig(file *ini.File) {
	mysql := file.Section("mysql")
	DbHost = mysql.Key("DbHost").String()
	DbPort = mysql.Key("DbPort").String()
	DbUser = mysql.Key("DbUser").String()
	DbPass = mysql.Key("DbPass").String()
	DbName = mysql.Key("DbName").String()
}

// 连接到MySQL
func connectMySQL(url string) {
	db, err := sql.Open("mysql", url)

	if err != nil {
		fmt.Println("Wrong configuration of [MySQL] in config file")
		panic(err)
	}

	// 最大连接数
	db.SetMaxOpenConns(100)
	// 闲置连接数
	db.SetMaxIdleConns(20)
	// 最大存活时间
	db.SetConnMaxLifetime(100 * time.Second)

	DB = db
	if err := DB.Ping(); err != nil {
		log.Printf("Connect to mysql server [%s] error\n", strings.Join([]string{DbHost, ":", DbPort}, ""))
		panic(err)
	} else {
		log.Printf("Connect to mysql server [%s] successfully\n", strings.Join([]string{DbHost, ":", DbPort}, ""))
	}
}
