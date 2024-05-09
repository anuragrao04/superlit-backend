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
	// the below is errr because err is already defined before. This error returns a function
	db.AutoMigrate(&models.User{}, &models.Classroom{}, &models.Assignment{}, &models.Submission{})
	return db, nil
}
