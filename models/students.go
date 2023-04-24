package models

import "gorm.io/gorm"

type Student struct {
	ID       uint    `gorm:"primary key;autoIncrement" json:"id"`
	FullName *string `json:"fullname"`
	Address  *string `json:"address"`
	Degree   *string `json:"degree"`
}

func MigrateStudents(db *gorm.DB) error {
	err := db.AutoMigrate(&Student{})
	return err
}
