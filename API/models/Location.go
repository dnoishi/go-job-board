package models

import "github.com/jinzhu/gorm"

type Location struct {
	gorm.Model
	LocationName string
}
