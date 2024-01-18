package models

import "gorm.io/gorm"

type User struct {
}

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}
