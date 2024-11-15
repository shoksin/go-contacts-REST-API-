package main

import (
	"fmt"
	"net/http"
	"os"

	app "github.com/shoksin/go-contacts-REST-API-/app"
	"github.com/shoksin/go-contacts-REST-API-/controllers"
	"github.com/shoksin/go-contacts-REST-API-/pkg/logging"

	"github.com/gorilla/mux"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("create router")
	router := mux.NewRouter()
	router.Use(app.JWTAuthentication)

	logger.Info("register user handlers")
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/api/contacts", controllers.GetContactsFor).Methods("GET")
	router.HandleFunc("/api/contacts/dalete/{name}", controllers.DeleteUserContact).Methods("DELETE")
	router.HandleFunc("/api/contacts/delete", controllers.DeleteAllUserContacts).Methods("DELETE")
	router.HandleFunc("/api/contacts/update/{contact_id}", controllers.UpdateUserContacts).Methods("PUT")
	router.HandleFunc("/api/contacts/update/{contact_id}", controllers.PatchUserContacts).Methods("PATCH")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
