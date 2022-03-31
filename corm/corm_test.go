package corm

import (
	"corms/models"
	"corms/setting"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)


func TestDoFunc(t *testing.T) {
	var d []models.CmpAreas
	err := models.DB().Select(&d, "select * from cmp_areas where id=100")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	for _, dd := range d {
		dd.Show()
	}

}

func TestParseTable(t *testing.T) {
	dev := models.CmpDevs{
		UserId:      4,

	}
	err := ParseTable(&dev)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestMain(m *testing.M) {
	path, _ := os.Getwd()
	config := fmt.Sprintf("%s%c..%cconf%capp_local.ini", path, filepath.Separator, filepath.Separator, filepath.Separator)
	setting.Setup(config)
	models.Setup()
	m.Run()
}