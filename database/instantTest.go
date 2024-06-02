package database

import (
	"errors"
	"github.com/anuragrao04/superlit-backend/models"
	"github.com/anuragrao04/superlit-backend/prettyPrint"
	"gorm.io/gorm"
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

	err = DB.Create(&newTest).Error
	prettyPrint.PrettyPrint(newTest)

	if err != nil {
		log.Println("Failed to create test: ", err)
		return "", "", err
	}
	return privateCode, publicCode, err
}

func GetInstantTest(publicCode string) (questions []models.Question, isActive bool, err error) {
	var test models.InstantTest
	if DB == nil {
		log.Fatal("Not Connected to Database")
	}
	err = DB.Preload("Questions").Preload("Questions.ExampleCases").Where("public_code = ?", publicCode).First(&test).Error
	if err != nil {
		log.Println("Failed to get test: ", err)
		return nil, false, err
	}

	if !test.IsActive {
		// test isn't active, don't allow the student to take it
		return nil, false, nil
	}

	// prettyPrint.PrettyPrint(test)
	return test.Questions, test.IsActive, nil
}

func ChangeActiveStatusInstantTest(active bool, privateCode string) (err error) {
	if DB == nil {
		log.Fatal("Not Connected to Database")
	}
	result := DB.Model(&models.InstantTest{}).Where("private_code = ?", privateCode).Update("is_active", active)
	if result.Error != nil {
		log.Println("Failed to change active status: ", err)
		return err
	}
	if result.RowsAffected == 0 {
		return errors.New("No test with that private code")
	}
	return nil
}

func UpsertSubmissionAndAnswers(instantTestID uint, universityID string, newAnswer models.Answer) error {
	// this creates a transaction. If any error occurs in any step, the entire transaction is rolled back
	// this ensures updating the submission and the answers is atomic

	return DB.Transaction(func(tx *gorm.DB) error {
		// Find existing submission
		var submission models.InstantTestSubmission
		err := tx.Where("instant_test_id = ? AND university_id = ?", instantTestID, universityID).First(&submission).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err // Error other than record not found
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Submission doesn't exist, create a new one
			submission = models.InstantTestSubmission{
				InstantTestID: instantTestID,
				UniversityID:  universityID,
				Answers:       []models.Answer{newAnswer}, // Start with the new answer
				TotalScore:    newAnswer.Score,
			}
			if err := tx.Create(&submission).Error; err != nil {
				return err
			}
		} else {
			// Submission exists. Now we look if an answer already exists for this question

			answerExists := false
			for i, existingAnswer := range submission.Answers {
				if existingAnswer.QuestionID == newAnswer.QuestionID {
					// Update existing answer

					submission.TotalScore -= existingAnswer.Score // Subtract old score
					submission.TotalScore += newAnswer.Score      // Add new score

					submission.Answers[i] = newAnswer // Replace the old answer
					answerExists = true
					// maybe in the future we'll be nice to the students and store the answer depending on which one scored higher
					// if the old score is better than new score, no updating
					break
				}
			}

			if !answerExists {
				// Answer doesn't exist, append the new answer
				submission.Answers = append(submission.Answers, newAnswer)
				submission.TotalScore += newAnswer.Score
			}

			if err := tx.Save(&submission).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
