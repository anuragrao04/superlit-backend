package main

import (
	"log"
	"net/http"

	"github.com/anuragrao04/superlit-backend/compile"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/run", compile.RunCode).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
