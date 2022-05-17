package models

import (
	"fmt"
	"time"
)

type MBase struct {
	ID uint64 `db:"id"`
}

type MCtime struct {
	MBase
	Ctime time.Time `db:"ctime"`
}

type MMtime struct {
	MCtime
	Mtime time.Time `db:"mtime"`
}


type TestDevs struct {
	MMtime
	Sn          string `db:"sn"`
	DevType     uint16 `db:"dev_type"`
	DevFunc     uint32 `db:"dev_func"`
	UserId      uint64 `db:"user_id"`
	VendorId    uint32 `db:"vendor_id"`
}

func (d *TestDevs) TableName() string {
	return "test_devs"
}

func (d *TestDevs) Show() {
	fmt.Printf("sn : %s, DevType: %d, UserId: %d, ctime: %d\n", d.Sn, d.DevType, d.UserId, d.Ctime.Unix())
}

