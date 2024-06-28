package controllers

import (
	"encoding/json"
	"fmt"
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var CreateContact = func(w http.ResponseWriter, r *http.Request) {

	// user, ok := r.Context().Value("user_id").(uint)
	// if !ok {
	// 	fmt.Println("User ============================", user)
	// 	u.Respond(w, u.Message(false, "Invalid user_id in the context"))
	// 	return
	// }

	contact := &models.Contact{}

	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
	}

	//contact.UserID = user
	resp := contact.Create()
	u.Respond(w, resp)

}

var GetContactsFor = func(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	fmt.Println(tokenHeader)
	if tokenHeader == "" {
		response := u.Message(false, "User is unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		u.Respond(w, response)
		return
	}
	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		response := u.Message(false, "Invalid auth token")
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		u.Respond(w, response)
		return
	}

	tokenPart := splitted[1]
	fmt.Println(tokenPart)
	tk := &models.Token{}

	token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})
	if err != nil {
		response := u.Message(false, "Token is not valid")
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		u.Respond(w, response)
		return
	}

	if !token.Valid {
		response := u.Message(false, "Token is not valid")
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		u.Respond(w, response)
		return
	}

	data := models.GetContacts(tk.UserId)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)

}
