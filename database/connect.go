package database

import (
	"github.com/anuragrao04/superlit-backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB // other packages can access this variable

func Connect(dbFilePath string) (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&models.User{}, &models.Classroom{}, &models.Assignment{}, &models.Submission{}, &models.TestCase{}, &models.Question{}, &models.InstantTest{})
	return DB, nil
}
