package database

import (
	// "errors"

	"errors"
	"log"

	"github.com/anuragrao04/superlit-backend/models"
	"gorm.io/gorm"
	// "github.com/anuragrao04/superlit-backend/prettyPrint"
)

func AddAssignmentToClassroom(assignment *models.Assignment, classroom *models.Classroom) error {
	DBLock.Lock()
	defer DBLock.Unlock()
	// result := DB.Model(classroom).Association("Assignments").Append(assignment)
	_ = DB.Model(classroom).Association("Assignments").Append(assignment)

	// FIXME: this result becomes null. I HAVE NO IDEA WHY.
	// Find out why
	// it always returns null. Even if it appends successfully.
	// Since result is null, I can't do a result.Error != nil. That would lead to a nil pointer dereference
	// So, I can't check for errors.
	// I have to assume that it worked.

	// prettyPrint.PrettyPrint(result)
	return nil
}

func GetAssignment(assignmentID uint) (*models.Assignment, error) {
	DBLock.Lock()
	defer DBLock.Unlock()
	var assignment models.Assignment
	result := DB.Preload("Questions.ExampleCases").Preload("Classrooms.Users").First(&assignment, assignmentID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &assignment, nil
}

// The difference between this function and
// GetAssignment is that this also fetches test cases, submissions, etc
func GetAssignmentForAIVerification(assignmentID uint) (models.Assignment, error) {
	var test models.Assignment
	DBLock.Lock()
	defer DBLock.Unlock()
	err := DB.Preload("Questions").Preload("Submissions").Preload("Submissions.Answers").First(&test, assignmentID).Error
	if err != nil {
		log.Println("Failed to get test: ", err)
		return models.Assignment{}, err
	}
	return test, nil
}

// the below function is used to update the submission and answers of a student
// if the submission exists, we update it. If not, it's created
func UpsertAssignmentSubmissionAndAnswers(assignmentID uint, userID uint, universityID string, newAnswer models.Answer) error {
	// this creates a transaction. If any error occurs in any step, the entire transaction is rolled back
	// this ensures updating the submission and the answers is atomic
	DBLock.Lock()
	defer DBLock.Unlock()
	return DB.Transaction(func(tx *gorm.DB) error {
		// Find existing submission
		var submission models.AssignmentSubmission
		err := tx.Preload("Answers").Preload("Answers.TestCases").Where("assignment_id = ? AND user_id = ?", assignmentID, userID).First(&submission).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err // Error other than record not found
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Submission doesn't exist, create a new one
			submission = models.AssignmentSubmission{
				AssignmentID: assignmentID,
				UserID:       userID,
				UniversityID: universityID,
				Answers:      []models.Answer{newAnswer}, // Start with the new answer
				TotalScore:   newAnswer.Score,
			}
			if err := tx.Create(&submission).Error; err != nil {
				return err
			}
		} else {
			// Submission exists. Now we look if an answer already exists for this question
			answerExists := false

			for _, existingAnswer := range submission.Answers {

				if existingAnswer.QuestionID == newAnswer.QuestionID {
					log.Println("Answer exists")
					// Update existing answer

					submission.TotalScore -= existingAnswer.Score // Subtract old score
					submission.TotalScore += newAnswer.Score      // Add new score
					// update the above in the database
					tx.Model(&submission).Where("id = ?", submission.ID).Update("total_score", submission.TotalScore)

					// Replace the old answer
					tx.Model(&submission).Association("Answers").Delete(existingAnswer)
					tx.Model(&submission).Association("Answers").Append(&newAnswer)

					// prettyPrint.PrettyPrint(submission)
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

func GetAssignmentSubmissions(assignmentID uint) (submissions []models.AssignmentSubmission, questionIDs []uint, err error) {
	var assignment models.Assignment

	DBLock.Lock()
	defer DBLock.Unlock()
	err = DB.Preload("Submissions").Preload("Questions").Preload("Submissions.Answers").First(&assignment, assignmentID).Error
	if err != nil {
		log.Println("Failed to get test: ", err)
		return nil, nil, err
	}

	for _, question := range assignment.Questions {
		questionIDs = append(questionIDs, question.ID)
	}

	return assignment.Submissions, questionIDs, nil
}

func SaveAssignment(assignment models.Assignment) error {
	err := DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&assignment).Error
	log.Println("Saved to db")
	return err
}
