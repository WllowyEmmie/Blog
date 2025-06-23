package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func InitializeDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=emituntun dbname=myblog port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	return db
}
func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Post{}, &Reaction{}, &Comment{})
}
