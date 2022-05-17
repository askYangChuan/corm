package corm

import (
	"fmt"
	"github.com/askYangc/corm/models"
	"github.com/askYangc/corm/setting"
	"os"
	"path/filepath"
	"testing"
	"time"
)


func TestCorm(t *testing.T) {
	tx := NewDB(models.DB())

	dev := models.CmpDevs{
		MtimeModel: models.MtimeModel{
			CtimeModel: models.CtimeModel{},
			Mtime:      time.Time{},
		},
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
	t.Log(dev)
}



func TestMain(m *testing.M) {
	path, _ := os.Getwd()
	config := fmt.Sprintf("%s%c..%cconf%capp_local.ini", path, filepath.Separator, filepath.Separator, filepath.Separator)
	setting.Setup(config)
	models.Setup()

	m.Run()
}