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

	router.POST("/run", compile.RunCode)

	s := &http.Server{
		Addr:         ":6969",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}
	log.Println("Listening on port 6969.")
	s.ListenAndServe()
}
