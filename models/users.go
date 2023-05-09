package models

import "gorm.io/gorm"

type User struct {
	ID       uint    `gorm:"primary key; autoIncrement" json:"id"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	return err
}
