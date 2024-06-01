package database

import (
	"github.com/anuragrao04/superlit-backend/models"
	"github.com/anuragrao04/superlit-backend/prettyPrint"
	"log"
)

func generateUniqueCodes() (privateCode, publicCode string) {
	for {
		privateCode := generateRandomCode()
		publicCode := generateRandomCode()
		var count int64
		// Check if the code already exists in the database
		// if it is , generate a new code, else return the code
		if DB == nil {
			log.Fatal("Not Connected to Database")
		}
		DB.Model(&models.InstantTest{}).Where("private_code = ? AND public_code = ?", privateCode, publicCode).Count(&count)
		if count == 0 {
			return privateCode, publicCode
		}
	}
}

func CreateInstantTest(questions []models.Question) (privateCode, publicCode string, err error) {
	// we assign a private and a public key to this test
	privateCode, publicCode = generateUniqueCodes()
	var newTest models.InstantTest
	newTest.PrivateCode = privateCode
	newTest.PublicCode = publicCode
	newTest.Questions = questions
	newTest.IsActive = true
	prettyPrint.PrettyPrint(newTest)

	err = DB.Create(&newTest).Error
	DB.Save(&newTest)

	if err != nil {
		log.Println("Failed to create test: ", err)
		return "", "", err
	}
	return privateCode, publicCode, err
}

func GetInstantTest(publicCode string) (questions []models.Question, err error) {
	var test models.InstantTest
	if DB == nil {
		log.Fatal("Not Connected to Database")
	}
	err = DB.Preload("Questions").Preload("Questions.ExampleCases").Where("public_code = ?", publicCode).First(&test).Error
	if err != nil {
		log.Println("Failed to get test: ", err)
		return nil, err
	}
	prettyPrint.PrettyPrint(test)
	return test.Questions, nil
}
