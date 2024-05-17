package database

import (
	"log"

	"github.com/anuragrao04/superlit-backend/models"
)

func CreateNewUser(UniversityID string, Name string, Email string, Password string, IsTeacher bool) (err error) {
	user := models.User{
		UniversityID: UniversityID,
		Name:         Name,
		Email:        Email,
		Password:     Password,
		IsTeacher:    IsTeacher,
		Classrooms:   []models.Classroom{},
	}
	err = DB.Create(&user).Error
	return err
}

func GetUserByUniversityIDPassword(UniversityID string, Password string) (user models.User, err error) {
	if DB == nil {
		log.Println("DB is nil. Not connected to database")
	}

	err = DB.Preload("Classrooms").Preload("Classrooms.Assignments").Where("university_id = ? AND password = ?", UniversityID, Password).First(&user).Error
	if err != nil {
		log.Println(err)
	}
	return user, err
}
