package controllers

import (
	"encoding/json"
	"go-contacts/models"
	"go-contacts/pkg/logging"
	u "go-contacts/utils"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func GetToken(w http.ResponseWriter, r *http.Request) *models.Token {
	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		logging.GetLogger().Error("Couldn't get token header.")
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		u.Respond(w, u.Message(false, "User is unauthorized"))
		return nil
	}

	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		logging.GetLogger().Error("The token does not consist of two parts.")
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		u.Respond(w, u.Message(false, "Invalid auth token"))
		return nil
	}

	tokedPart := splitted[1]
	tk := &models.Token{}

	token, err := jwt.ParseWithClaims(tokedPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})
	if err != nil {
		logging.GetLogger().Error("Couldn't parse JWT.")
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		u.Respond(w, u.Message(false, "Token is not valid!"))
		return nil
	}

	if !token.Valid {
		logging.GetLogger().Info("Token isn't valid")
		response := u.Message(false, "Token is not valid")
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		u.Respond(w, response)
		return nil
	}

	return tk
}

var CreateContact = func(w http.ResponseWriter, r *http.Request) {
	contact := &models.Contact{}
	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		logging.GetLogger().Errorf("the request body could not be decoded")
		u.Respond(w, u.Message(false, "Error while decoding request body"))
	}
	contact.UserID = GetToken(w, r).UserId
	resp := contact.Create()
	u.Respond(w, resp)
}

var GetContactsFor = func(w http.ResponseWriter, r *http.Request) {
	data := models.GetContacts(GetToken(w, r).UserId)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)

}

var DeleteUserContact = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nameStr := vars["name"]
	data := models.DeleteContact(GetToken(w, r).UserId, nameStr)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)

}

var DeleteAllUserContacts = func(w http.ResponseWriter, r *http.Request) {
	data := models.DeleteContacts(GetToken(w, r).UserId)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var UpdateUserContacts = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contactIDString := vars["contact_id"]
	contactID, err := strconv.Atoi(contactIDString)
	if err != nil {
		logging.GetLogger().Error("Couldn't convert concactID(string) to integer")
		u.Respond(w, u.Message(false, "Wrong contactID"))
	}

	contact := &models.Contact{}
	err = json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		logging.GetLogger().Error("the request body could not be decoded")
		u.Respond(w, u.Message(false, "Error while decoding request body"))
	}
	data := models.UpdateContact(GetToken(w, r).UserId, uint(contactID), contact)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}

var PatchUserContacts = func(w http.ResponseWriter, r *http.Request) {
	contactIDString := mux.Vars(r)["contact_id"]
	contactID, err := strconv.Atoi(contactIDString)
	if err != nil {
		logging.GetLogger().Error("Couldn't convert concactID(string) to integer")
		u.Respond(w, u.Message(false, "Wrong contactID"))
	}

	contact := &models.Contact{}
	err = json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		logging.GetLogger().Error("the request body could not be decoded")
		u.Respond(w, u.Message(false, "Error while decoding request body:"))
	}

	data := models.PatchContact(GetToken(w, r).UserId, uint(contactID), contact)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
