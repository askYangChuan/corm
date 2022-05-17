package main

import (
	"fmt"
	"github.com/askYangc/corm/corm"
	"github.com/askYangc/corm/models"
	"github.com/askYangc/corm/setting"
)

func main() {
	setting.Setup("conf/app.ini")
	models.Setup()
	db := corm.NewDB(models.DB())

	dev := models.CmpDevs{}

	models.DB().Get(&dev, "select * from cmp_devs where id=5 limit 1")

	dev.DevType = 500
	dev.Sn = "update"
	db.InsertOrUpdate(&dev)

	fmt.Println(dev.MtimeModel)
	fmt.Println(dev)
	//db.Insert(xxx)
	//db.Select(xxx).Order().Limit()
	//db.update(xxx)
}