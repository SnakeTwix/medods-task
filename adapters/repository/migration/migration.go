package migration

import (
	"gorm.io/gorm"
	"medods-api/adapters/repository/model"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&model.Token{})

	return err
}
