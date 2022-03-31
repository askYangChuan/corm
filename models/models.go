package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type BaseModel struct {
	ID uint64 `db:"id"`
}

type CtimeModel struct {
	BaseModel
	Ctime time.Time `db:"ctime"`
}

type MtimeModel struct {
	CtimeModel
	Mtime time.Time `db:"mtime"`
}


