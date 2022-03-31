package main

import (
	"corms/corm"
	"corms/models"
)

func main() {
	models.Setup()

	main_test()
}

func main_test() {
	db := corm.NewDB(models.DB())

	/*
	corm.ParseTable(&models.CmpDevs{})
	corm.Show(&models.CmpDevs{})

	 */
	//插入
	db.Create(&models.CmpDevs{})
	//更新
	//db.Update(&models.CmpDevs{})
/*
	db.Select("Name", "Age", "CreatedAt").Create(&user)
   // INSERT INTO `users` (`name`,`age`,`created_at`) VALUES ("jinzhu", 18, "2020-07-04 11:05:21.775")
 */
	//db.Where("dev_type=?", 10).Get(&models.CmpDevs{})

}
