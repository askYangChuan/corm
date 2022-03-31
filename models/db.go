package models

import (
	"corms/setting"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB


func mysqlInit() {
	var err error
	db, err = sqlx.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("connect ok")
		panic(err)
	}

	//最大空闲连接数
	//db.SetMaxIdleConns(5)
	//最大连接数
	//db.SetMaxOpenConns(10)
	//自动重连尝试目前是2次
}

func dbInit() {
	switch setting.DatabaseSetting.Type {
	case "mysql":
		mysqlInit()
	}
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	db.Close()
}

func DB() *sqlx.DB {
	return db
}

// Setup initializes the database instance
func Setup() {
	dbInit()
}