package corm

import (
	"fmt"
	"github.com/askYangc/corm/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"testing"
)

var testDb *sqlx.DB

type Database struct {
	Type         string
	User         string
	Password     string
	Host         string
	Name         string
}

func TestCorm(t *testing.T) {
	tx := NewDB(testDb)

	dev := models.TestDevs{
		Sn:         "haha",
		DevType:    1,
		DevFunc:    2,
		UserId:     3,
		VendorId:   4,
	}

	x := tx.Insert(&dev)
	if x.Error != nil {
		t.Log(x.Error)
		t.Fail()
	}

	dev.DevType = 100
	x = tx.Update(&dev)
	if x.Error != nil {
		t.Log(x.Error)
		t.Fail()
	}

	dev.DevFunc = 200
	x = tx.InsertOrUpdate(&dev)
	if x.Error != nil {
		t.Log(x.Error)
		t.Fail()
	}

	x = tx.Delete(&dev)
	if x.Error != nil {
		t.Log(x.Error)
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	tx := NewDB(testDb)

	dev := models.TestDevs{}
	x := tx.Get(&dev, "dev_type=? and sn=?", 4, "1")
	if x.Error != nil {
		t.Log(x.Error)
		t.Fail()
	}
	fmt.Println(dev)
}

func mysqlInit(databaseSetting *Database) {
	var err error
	testDb, err = sqlx.Open(databaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		databaseSetting.User,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.Name))
	if err != nil {
		panic(err)
	}
	err = testDb.Ping()
	if err != nil {
		fmt.Println("connect failed")
		panic(err)
	}
	fmt.Println("connect ok")
	//最大空闲连接数
	//db.SetMaxIdleConns(5)
	//最大连接数
	//db.SetMaxOpenConns(10)
	//自动重连尝试目前是2次
}

func TestMain(m *testing.M) {
	databaseSetting := Database{
		Type:     "mysql",
		User:     "root",
		Password: "123456",
		Host:     "10.135.255.202:3306",
		Name:     "test",
	}
	mysqlInit(&databaseSetting)
	m.Run()
	testDb.Close()
}