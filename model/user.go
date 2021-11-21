package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null" json:"name"`
	Telephone string `gorm:"varchar(11);not null;unique" json:"telephone"`
	Password  string `gorm:"size:255;not null" json:"password"`
}
