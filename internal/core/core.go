package core

import "gorm.io/gorm"

type Core struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Core {
	return &Core{
		db: db,
	}
}
