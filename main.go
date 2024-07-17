package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/anuragrao04/superlit-backend/AI"
	"github.com/anuragrao04/superlit-backend/assignments"
	"github.com/anuragrao04/superlit-backend/auth"
	"github.com/anuragrao04/superlit-backend/classroom"
	"github.com/anuragrao04/superlit-backend/compile"
	"github.com/anuragrao04/superlit-backend/database"
	"github.com/anuragrao04/superlit-backend/googleSheets"
	"github.com/anuragrao04/superlit-backend/instantTest"
	"github.com/anuragrao04/superlit-backend/mailers"
	"github.com/anuragrao04/superlit-backend/tokens"
)

func main() {
	router := gin.Default()
	// initialise the router

	// load envs
	godotenv.Load()

	// load the JWT Secret Key
	err := tokens.LoadPrivateKey()
	if err != nil {
		log.Println(err)
		log.Fatal("Error loading private key")
	} else {
		log.Println("Loaded JWT private key")
	}

	// connect to the database
	_, err = database.Connect("holy.db")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to the database.")
	}

	// connect to google sheets api
	err = googleSheets.Connect()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to Google Sheets API.")
	}

	// connect to the email server
	err = mailers.Connect()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to the email server.")
	}

	// These are the routes available. No other routes apart from these will be available. All routes must be defined here.

	// this route is for running the code. Only run. No scoring no judging
	router.POST("/run", compile.RunCode)

	// this route is for creating a new user
	router.POST("/auth/signup", auth.SignUp)

	// this route is for signing in with a University ID
	// TODO signing with email
	router.POST("/auth/signinwithuniversityid", auth.SignInWithUniversityID)

	router.GET("/auth/isteacher", tokens.VerifyToken, auth.IsTeacherFromToken)

	// this route is for users who've forgotten their password.
	// It takes in University ID and sends a password reset link to their email
	router.POST("/auth/forgotpassword", auth.ForgotPassword)

	// this route is for users who have acquired their jwt token to reset the password
	// It takes in the jwt token and new password and changes the password in the database
	router.POST("/auth/resetpassword", auth.ResetPassword)

	// this route is used once signed in
	// to get basic user details from token
	router.GET("/auth/getuser", tokens.VerifyToken, auth.GetUserFromToken)

	// this route is for creating a new classroom
	router.POST("/classroom/create", tokens.VerifyToken, classroom.CreateClassroom)

	// add a user to a classroom. Regardless of if a user is a teacher or a student, they hit this route
	router.POST("/classroom/adduser", tokens.VerifyToken, classroom.AddUserToClassroom)

	// list assignments of a particular classroom
	router.POST("/assignment/listassignments", tokens.VerifyToken, classroom.ListAssignments)

	// create an assignment
	router.POST("/assignment/createassignment", tokens.VerifyToken, assignments.CreateAssignment)

	// get assignment data. Used when a student attempts an assignment. It returns questions and stuff
	router.POST("/assignment/get", tokens.VerifyToken, assignments.GetAssignment)

	// submit an assignment. Does scoring and stores in DB
	router.POST("/assignment/submit", tokens.VerifyToken, assignments.Submit)

	// get submissions for an assignment
	router.POST("/assignment/getsubmissions", tokens.VerifyToken, assignments.GetAssignmentSubmissions)

	// get submissions for a student
	router.POST("/assignment/getsubmissionforstudent", tokens.VerifyToken, assignments.GetStudentSubmission)

	// google sheet population for assignment
	router.POST("/assignment/populategooglesheet", tokens.VerifyToken, googleSheets.PopulateGoogleSheetAssignment)

	// AI Verification of assignment
	router.POST("/assignment/aiverify", AI.AIVerifyConstraintsAssignment)

	// This route is used when a teacher wants to edit an assignment
	// It will send the existing assignment data to the teacher
	router.POST("/assignment/getforedit", tokens.VerifyToken, assignments.GetAssignmentForEdit)

	// The below route is used to save the edited assignment
	// obtained from the above route
	router.POST("/assignment/saveedited", tokens.VerifyToken, assignments.SaveEditedAssignment)

	// INSTANT TEST STUFF
	// Create an instant test
	router.POST("/instanttest/create", instantTest.CreateTest)

	// get test data. Used when starting a test
	router.POST("/instanttest/get", instantTest.GetInstantTest)

	// change active status of a test
	router.POST("/instanttest/changeactive", instantTest.ChangeActive)

	// submit a question from the instant test
	router.POST("/instanttest/submit", instantTest.Submit)

	// get submissions for an instant test
	router.POST("/instanttest/getsubmissions", instantTest.GetSubmissions)

	// AI Verification of instant test
	router.POST("/instanttest/aiverify", AI.AIVerifyConstraintsInstantTest)

	// Populate google sheet with instant test submission data
	router.POST("/instanttest/populategooglesheet", googleSheets.PopulateGoogleSheetInstantTest)

	s := &http.Server{
		Addr:         ":6969",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Listening on port 6969.")
	s.ListenAndServe()
}
