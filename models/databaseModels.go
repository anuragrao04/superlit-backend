package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UniversityID string      `json:"universityID"` // this may be employee ID or student ID (EMP ID and SRN in terms of PES)
	Name         string      `json:"name"`
	Email        string      `json:"email"`
	Password     string      `json:"password"`
	IsTeacher    bool        `json:"isTeacher"`
	Classrooms   []Classroom `gorm:"many2many:user_classroom;" json:"classrooms"`
}

// aka class
type Classroom struct {
	gorm.Model
	Name        string       `json:"name"`
	Code        string       `json:"code"` // this is the code that students will use to join the classroom
	Users       []User       `json:"users" gorm:"many2many:user_classroom;"`
	Assignments []Assignment `json:"assignments" gorm:"many2many:assignment_classroom;"`
}

// aka test
type Assignment struct {
	gorm.Model
	Name        string       `json:"name"`
	Description string       `json:"description"`
	StartTime   time.Time    `json:"startTime"`
	EndTime     time.Time    `json:"endTime"`
	Classrooms  []Classroom  `gorm:"many2many:assignment_classroom;" json:"classrooms"`
	Questions   []Question   `gorm:"foreignKey:AssignmentID" json:"questions"`
	Submissions []Submission `gorm:"foreignKey:AssignmentID" json:"submissions"`

	// this is the classrooms in which the assignment is assigned.
	// Since every classroom can have multiple assignments, and one assignment may be assigned to multiple classrooms, we have a many to many relationship
}

// Instant tests are tests that are created on the fly by teachers.
// Private code is not shared. This code is used to edit the test
// Public code is shared with students. This code is used to attempt the test
type InstantTest struct {
	gorm.Model
	Email       string                  `json:"email"`
	PrivateCode string                  `json:"privateCode"`
	PublicCode  string                  `json:"publicCode"`
	IsActive    bool                    `json:"isActive"`
	Questions   []Question              `gorm:"foreignKey:InstantTestID" json:"questions"`
	Submissions []InstantTestSubmission `gorm:"foreignKey:InstantTestID" json:"submissions"`
}

type InstantTestSubmission struct {
	gorm.Model
	InstantTestID uint     `json:"instantTestID"`
	UniversityID  string   `json:"universityID"`
	Answers       []Answer `json:"answers"`
	TotalScore    int      `json:"totalScore"`
}

type Submission struct {
	gorm.Model
	AssignmentID uint     `json:"assignmentID"`
	UserID       uint     `json:"userID"`
	User         User     `json:"user"`
	Answers      []Answer `json:"answers"`
	TotalScore   int      `json:"totalScore"`
}

type Answer struct {
	gorm.Model
	SubmissionID            uint               `json:"submissionID"`
	InstantTestSubmissionID uint               `json:"instantTestSubmissionID"`
	QuestionID              uint               `json:"questionID"`
	Code                    string             `json:"code"`
	TestCases               []VerifiedTestCase `json:"testCases"`
	Score                   int                `json:"score"`      // total score for this particular question
	AIVerified              bool               `json:"AIVerified"` // if AI has verified the code
	AIVerdict               bool               `json:"AIVerdict"`  // if AI has verified the code, this is the verdict. If true, it means it's aproved. else something is fishy
}

// The below gymnastics is because GORM doesn't support storing []string directly.
// So we have to implement the Scanner and Valuer interfaces to convert the []string to a varchar and back

type Question struct {
	gorm.Model
	AssignmentID   uint              `json:"assignmentID"`
	InstantTestID  uint              `json:"instantTestID"`
	Title          string            `json:"title"`
	Question       string            `json:"question"`
	ExampleCases   []ExampleTestCase `json:"exampleCases"`
	PreWrittenCode string            `json:"preWrittenCode"`
	TestCases      []TestCase        `json:"testCases"`
	Constraints    pq.StringArray    `json:"constraints" gorm:"type:text[]"` // list of constraints that the code must satisfy. Used for AI Verification
	Languages      pq.StringArray    `json:"languages" gorm:"type:text[]"`   // languages the student is allowed to submit in
}

type ExampleTestCase struct {
	gorm.Model
	QuestionID     uint   `json:"questionID"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expectedOutput"`
	Score          int    `json:"score"`
	Explanation    string `json:"explanation"`
}

type TestCase struct {
	gorm.Model
	QuestionID     uint   `json:"questionID"`
	Input          string `'json:"input"`
	ExpectedOutput string `json:"expectedOutput"`
	Score          int    `json:"score"`
}

type VerifiedTestCase struct {
	AnswerID       uint   `json:"answerID"`
	Passed         bool   `json:"passed"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expectedOutput"`
	ProducedOutput string `json:"producedOutput"`
}
