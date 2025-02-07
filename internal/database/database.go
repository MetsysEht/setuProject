package database

import (
	"github.com/MetsysEht/setuProject/internal/kycVerification/model"
	"github.com/MetsysEht/setuProject/pkg/gormDatabase"
	"github.com/MetsysEht/setuProject/pkg/logger"
	"gorm.io/gorm"
)

func GetDatabase(c *gormDatabase.Config) (*gorm.DB, error) {
	db, err := gormDatabase.CreateGormDatabase(c)
	if err != nil {
		logger.L.Fatal("Could not connect to DB")
		return nil, err
	}
	updateTables(db)
	return db, nil
}

func updateTables(db *gorm.DB) {
	_ = db.AutoMigrate(&model.PANVerification{})
	_ = db.AutoMigrate(&model.RPDVerification{})
}
