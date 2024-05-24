package main

import (
	"log"
	"net/http"
	"time"

	"github.com/anuragrao04/superlit-backend/auth"
	"github.com/anuragrao04/superlit-backend/classRoom"
	"github.com/anuragrao04/superlit-backend/compile"
	"github.com/anuragrao04/superlit-backend/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
	// initialise the router

	// load envs
	godotenv.Load()

	// connect to the database
	_, err := database.Connect("holy.db")
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
	router.POST("./auth/signinwithuniversityid", auth.SignInWithUniversityID)

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
