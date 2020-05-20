package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"web_server/dataBase"
)

func main() {
	router := mux.NewRouter()
	db, err := dataBase.NewDataBase("postgres", "admin", "meet_and_go", "35.228.190.27", 5432)
	if err != nil {
		log.Fatalf("DataBase error")
	}
	handlers := NewHandlers(db)
	router.HandleFunc("/registration", handlers.Registration).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/profiles/{id}", handlers.GetUserInfo).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
