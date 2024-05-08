package main

import (
	"log"
	"net/http"
	"time"

	"github.com/anuragrao04/superlit-backend/compile"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// initialise the router

	// These are the routes available. No other routes apart from these will be available. All routes must be defined here.
	router.POST("/run", compile.RunCode)

	s := &http.Server{
		Addr:         ":6969",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Listening on port 6969.")
	s.ListenAndServe()
}
