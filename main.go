package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"web_server/dataBase"
)

func main() {
	config, err := getConfig("config.json")
	if err != nil {
		log.Fatalf(err.Error())
	}
	db, err := dataBase.NewDataBase(config[dataBaseField].(*dataBase.DBConfig))
	if err != nil {
		log.Fatalf("DataBase error")
	}

	router := mux.NewRouter()
	handlers := NewHandlers(db, config[logField].(bool))
	router.HandleFunc("/registration", handlers.Registration).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/profiles/{id}", handlers.GetUserInfo).Methods("GET")
	router.HandleFunc("/search", handlers.SearchUser).Methods("GET")

	//For developers
	router.HandleFunc("/developers/getAccount", handlers.GetDeveloperAccount).Methods("GET")
	router.HandleFunc("/developers/postTag", handlers.PostTag).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
