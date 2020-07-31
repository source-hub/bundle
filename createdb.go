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
