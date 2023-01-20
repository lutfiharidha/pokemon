package db

import (
	"github.com/lutfiharidha/pokemon/app/types"
	"gorm.io/gorm"
)

func Migrator(db *gorm.DB) {
	db.AutoMigrate(&types.Log{})
}
