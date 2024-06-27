package controllers

import (
	"encoding/json"
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		//Переданный параметр пути не является целым числом
		u.Respond(w, u.Message(false, "There was an error in your request"))
	}

	data := models.GetContacts(uint(id))
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)

}
