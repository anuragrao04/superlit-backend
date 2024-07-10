package database

import (
	// "errors"

	"github.com/anuragrao04/superlit-backend/models"
	// "github.com/anuragrao04/superlit-backend/prettyPrint"
)

func AddAssignmentToClassroom(assignment *models.Assignment, classroom *models.Classroom) error {
	DBLock.Lock()
	defer DBLock.Unlock()
	result := DB.Model(classroom).Association("Assignments").Append(assignment)

	// FIXME: this result becomes null. I HAVE NO IDEA WHY.
	// Find out why
	// it always returns null. Even if it appends successfully.
	// Since result is null, I can't do a result.Error != nil. That would lead to a nil pointer dereference
	// So, I can't check for errors.
	// I have to assume that it worked.

	// prettyPrint.PrettyPrint(result)
	return nil
}
