package database

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/anuragrao04/superlit-backend/models"
)

func CreateClassroom(Name string, TeacherID uint) (*models.Classroom, error) {
	// first we make sure given TeacherID is a teacher
	var teacher models.User

	DBLock.Lock()
	defer DBLock.Unlock()
	err := DB.First(&teacher, TeacherID).Error
	if err != nil {
		return nil, errors.New("teacher not found")
	}
	if !teacher.IsTeacher {
		return nil, errors.New("user is not a teacher")
	}
	code := generateUniqueCode()
	classroom := models.Classroom{
		Name: Name,
		Code: code,
	}

	err = DB.Create(&classroom).Error
	if err != nil {
		return nil, err
	}

	// Add the teacher to the classroom
	err = AddUserToClassroom(TeacherID, classroom.Code, true)
	return &classroom, err
}

func AddUserToClassroom(UserID uint, ClassroomCode string, IsTeacher bool) error {
	var user models.User
	var classroom models.Classroom

	DBLock.Lock()
	defer DBLock.Unlock()

	if err := DB.First(&user, UserID).Error; err != nil {
		return errors.New("user not found")
	}

	if err := DB.Where("code = ?", ClassroomCode).First(&classroom).Error; err != nil {
		return errors.New("classroom not found")
	}

	classroom.Users = append(classroom.Users, user) // Add the user to the classroom
	// updating it here will add the association in the user struct as well. GORM takes care of it

	if err := DB.Save(&classroom).Error; err != nil {
		return fmt.Errorf("failed to associate user with classroom: %w", err)
	}
	return nil
}

// this function is only used by the CreateClassroom function
// The database is already locked there. So we shouldn't attempt to acquire the mutex here
func generateUniqueCode() string {
	for {
		code := generateRandomCode()
		var count int64
		// Check if the code already exists in the database
		// if it is , generate a new code, else return the code
		DB.Model(&models.Classroom{}).Where("code = ?", code).Count(&count)
		if count == 0 {
			return code
		}
	}
}

func generateRandomCode() string {
	code := rand.Intn(900000) + 100000 // Generates a number between 100000 and 999999
	return strconv.Itoa(code)
}
