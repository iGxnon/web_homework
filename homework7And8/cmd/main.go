package main

import "fmt"

func main() {
	//dao.InitializeDefault()
	//api.RegisterRouter()
	fmt.Println("CREATE TABLE `tree_hollow` (" +
		"`id` BIGINT(20) NOT NULL AUTO_INCREMENT," +
		"`prefix` VARCHAR(20) DEFAULT ''," +
		"`model_json` MEDIUMTEXT NOT NULL," +
		"`model_serialized_img` MEDIUMTEXT NOT NULL," +
		"`loc` NOT NULL TINYTEXT," +
		"`form_title` NOT NULL VARCHAR(20) DEFAULT ''," +
		"`form_texts` NOT NULL TINYTEXT," +
		"`animation` TINYTEXT," +
		"`particle` TINYTEXT," +
		"PRIMARY KEY(`id`, `prefix`)" + // 	双主键
		")ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;")
}
