package main

import (
	"fmt"
	app "go-contacts/app"
	"go-contacts/controllers"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Use(app.JWTAuthentication) //добавляем middleware проверки JWT-токена

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/api/user/{id}/contacts", controllers.GetContactsFor).Methods("GET") //  user/2/contacts

	router.Use(app.JWTAuthentication)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
