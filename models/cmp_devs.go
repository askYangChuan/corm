package models

import "fmt"

type CmpDevs struct {
	MtimeModel
	Sn          string `db:"sn"`           //color_system_id ,我们分配的
	DevType     uint16 `db:"dev_type"`     //设备类型	1,染色网关，2信令网关
	DevFunc     uint32 `db:"dev_func"`     //设备功能，BIT(0) 染色功能; BIT(1) 信令功能
	UserId      uint64 `db:"user_id"`      //设备所属安全管理员用户ID
	VendorId    uint32 `db:"vendor_id"`    //厂商ID
}

func (d *CmpDevs) TableName() string {
	return "cmp_devs"
}

func (d *CmpDevs) Show() {
	fmt.Printf("sn : %s, DevType: %d, UserId: %d, ctime: %d\n", d.Sn, d.DevType, d.UserId, d.Ctime.Unix())
}

//区域表
type CmpAreas struct {
	CtimeModel
	AreaId   uint32 `db:"area_id"`
	AreaName string `db:"area_name"`
}

func (d *CmpAreas) Show() {
	fmt.Printf("AreaId : %d, AreaName: %s, ctime: %d\n", d.AreaId, d.AreaName, d.Ctime.Unix())
}