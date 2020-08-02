package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"bundle/models"
)

func Migrate(db *gorm.DB) {

	//
	println("Performing the database migration")
	db.AutoMigrate(
		&models.User{},
		&models.Catalog{},
		&models.Item{},
		&models.Item_Catalog{},
	)
	db.Model(&models.Catalog{}).AddForeignKey("user_id","users(id)","cascade","cascade")
	db.Model(&models.Item_Catalog{}).AddUniqueIndex("idx_catalog_id_item_id","catalog_id","item_id")
	println("done performing migrations")
}

func Drop(db *gorm.DB) {
	db.DropTableIfExists(
		&models.User{},
		&models.Catalog{},
		&models.Item{},
		&models.Item_Catalog{},
	)
}
