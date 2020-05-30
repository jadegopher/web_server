package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"web_server/dataBase"
)

func main() {
	config, err := getConfig("secureConfig.json")
	if err != nil {
		log.Fatalf(err.Error())
	}
	db, err := dataBase.NewDataBase(config[dataBaseField].(*dataBase.DBConfig))
	if err != nil {
		log.Fatalf(err.Error())
	}

	router := mux.NewRouter()
	handlers := NewHandlers(db, config[logField].(bool))
	router.HandleFunc("/registration", handlers.Registration).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/profiles/{id}", handlers.GetUserInfo).Methods("GET")
	router.HandleFunc("/search", handlers.SearchUser).Methods("GET")
	router.HandleFunc("/delete", handlers.DeleteUser).Methods("POST")
	router.HandleFunc("/tags", handlers.GetTags).Methods("GET")
	router.HandleFunc("/tags/add", handlers.AddTagsToUser).Methods("POST")
	router.HandleFunc("/tags/get/{id}", handlers.GetUserTags).Methods("GET")
	router.HandleFunc("/tasks/{taskName}", handlers.GetTaskInfo).Methods("GET")
	router.HandleFunc("/tasks/{taskName}/tags", handlers.GetTaskTags).Methods("GET")
	//For developers
	router.HandleFunc("/developers/getAccount", handlers.GetDeveloperAccount).Methods("GET")
	router.HandleFunc("/developers/postTag", handlers.PostTag).Methods("POST")
	router.HandleFunc("/developers/tasks/post", handlers.PostTask).Methods("POST")
	router.HandleFunc("/developers/tasks/addTags/{taskName}", handlers.AddTagsToTask).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
