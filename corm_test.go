package corm

import (
	"fmt"
	"github.com/askYangc/corm/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"testing"
	"time"
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
	x := tx.Get(&dev, "dev_type=? limit 1", 4)
	if x.Error != nil {
		t.Log(x.Error)
		t.Fail()
	}
	fmt.Println(dev)
}

func TestSelect(t *testing.T) {
	tx := NewDB(testDb)

	var devs []models.TestDevs
	x := tx.Select(&devs, "dev_type=1 limit ?,?", 1, 2)
	if x.Error != nil {
		t.Log(x.Error)
		t.Fail()
	}

	fmt.Println(devs)
}

func TestInsert(t *testing.T) {
	tx := NewDB(testDb)

	dev := models.TestDevs{
		Sn:         "haha1",
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
}

func TestTransTheOne(t *testing.T) {
	c := NewDB(testDb)

	tx, err := c.Beginx()
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	defer tx.Rollback()

	dev := models.TestDevs{
		Sn:       "tran1",
		DevType:  6,
		DevFunc:  6,
		UserId:   6,
		VendorId: 6,
	}

	tx.Insert(&dev)

	tx.Commit()
}

func TestTransTheTwo(t *testing.T) {
	c := NewDB(testDb)

	originTx, err := testDb.Beginx()
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	tx := c.SetTx(originTx)

	defer tx.Rollback()

	dev := models.TestDevs{
		Sn:       "tran2",
		DevType:  6,
		DevFunc:  6,
		UserId:   6,
		VendorId: 6,
	}

	tx.Insert(&dev)
	originTx.Commit()
}

func TestTransLock(t *testing.T) {
	c := NewDB(testDb)

	tx, err := c.Beginx()
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	defer tx.Rollback()

	var dev models.TestDevs
	tx.Get(&dev, "sn='wo' limit 1 for update")

	dev.DevType = 9

	fmt.Println(dev)

	tx.Update(&dev)

	time.Sleep(time.Second)

	tx.Commit()
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
		Password: "",
		Host:     "127.0.0.1:3306",
		Name:     "test",
	}

	mysqlInit(&databaseSetting)

	logger, _ := zap.NewDevelopment()
	SetLogger(logger)

	m.Run()
	testDb.Close()
}