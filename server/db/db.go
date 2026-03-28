package db

import (
	"fmps/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init(dbPath string) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = database.AutoMigrate(
		&models.Setting{},
		&models.Member{},
		&models.Rule{},
		&models.Record{},
	)
	if err != nil {
		return nil, err
	}

	return database, nil
}
