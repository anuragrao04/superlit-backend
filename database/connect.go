package database

import (
	"github.com/anuragrao04/superlit-backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(dbFilePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{}, &models.Classroom{}, &models.Assignment{}, &models.Submission{}, &models.TestCase{}, &models.Question{})
	return db, nil
}
