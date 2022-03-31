package gorm_models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func test() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	db.Select("").Find(&db)
	db.Select("Name", "Age", "CreatedAt").Create(&user)
}