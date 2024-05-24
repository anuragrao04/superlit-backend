package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UniversityID string // this may be employee ID or student ID (EMP ID and SRN in terms of PES)
	Name         string
	Email        string
	Password     string
	IsTeacher    bool
	Classrooms   []Classroom `gorm:"many2many:user_classroom;"`
}

// aka class
type Classroom struct {
	gorm.Model
	Name        string
	Code        string       // this is the code that students will use to join the classroom
	Users       []User       `gorm:"many2many:user_classroom;"`
	Assignments []Assignment `gorm:"many2many:assignment_classroom;"`
}

// aka test
type Assignment struct {
	gorm.Model
	Name        string
	Description string
	StartTime   *time.Time
	EndTime     *time.Time
	Classrooms  []Classroom `gorm:"many2many:assignment_classroom;"`
	Questions   []Question
	Submissions []Submission
	// this is the classrooms in which the assignment is assigned.
	// Since every classroom can have multiple assignments, and one assignment may be assigned to multiple classrooms, we have a many to many relationship
}

type Submission struct {
	gorm.Model
	AssignmentID uint
	UserID       uint
	User         User
	Answers      []Answer
	TotalScore   int
}

type Answer struct {
	gorm.Model
	SubmissionID    uint
	QuestionID      uint
	Question        Question
	Code            string
	TestCasesPassed []TestCase `gorm:"many2many:answer_testcases_passed;"`
	Score           int        // total score for this particular question
}

type Question struct {
	gorm.Model
	AssignmentID   uint
	Title          string
	Question       string
	ExampleCases   []TestCase
	PreWrittenCode string
	TestCases      []TestCase
}

type TestCase struct {
	gorm.Model
	QuestionID     uint
	Input          string
	ExpectedOutput string
	Score          int // score achieved by the student if the output matches the expected output
}
