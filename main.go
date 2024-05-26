package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/anuragrao04/superlit-backend/auth"
	"github.com/anuragrao04/superlit-backend/classroom"
	"github.com/anuragrao04/superlit-backend/compile"
	"github.com/anuragrao04/superlit-backend/database"
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

	// These are the routes available. No other routes apart from these will be available. All routes must be defined here.

	// this route is for running the code. Only run. No scoring no judging
	router.POST("/run", compile.RunCode)

	// this route is for creating a new user
	router.POST("/auth/signup", auth.SignUp)

	// this route is for signing in with a University ID
	// TODO signing with email
	router.POST("/auth/signinwithuniversityid", auth.SignInWithUniversityID)

	// this route is for users who've forgotten their password.
	// It takes in University ID and sends a password reset link to their email
	router.POST("/auth/forgotpassword", auth.ForgotPassword)

	// this route is for users who have acquired their jwt token to reset the password
	// It takes in the jwt token and new password and changes the password in the database
	router.POST("/auth/resetpassword", auth.ResetPassword)

	// this route is for creating a new classroom
	router.POST("/classroom/create", classroom.CreateClassroom)

	// add a user to a classroom. Regardless of if a user is a teacher or a student, they hit this route
	router.POST("/classroom/adduser", classroom.AddUserToClassroom)

	s := &http.Server{
		Addr:         ":6969",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Listening on port 6969.")
	s.ListenAndServe()
}
