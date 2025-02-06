package database

import (
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
	return db, nil
}
