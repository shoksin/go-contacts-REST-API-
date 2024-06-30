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
	router.HandleFunc("/api/contacts", controllers.GetContactsFor).Methods("GET")
	router.HandleFunc("/api/contacts/dalete/{name}", controllers.DeleteUserContact).Methods("DELETE")
	router.HandleFunc("/api/contacts/delete", controllers.DeleteAllUserContacts).Methods("DELETE")
	router.HandleFunc("/api/contacts/update/{contact_id}", controllers.UpdateUserContacts).Methods("PUT")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
