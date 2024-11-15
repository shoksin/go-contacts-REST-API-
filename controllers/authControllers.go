package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/shoksin/go-contacts-REST-API-/models"
	"github.com/shoksin/go-contacts-REST-API-/pkg/logging"
	u "github.com/shoksin/go-contacts-REST-API-/utils"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //декодирует тело запроса в struct и завершается неудачно в случае ошибки
	if err != nil {
		logging.GetLogger().Errorf("the request body could not be decoded")
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create()
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		logging.GetLogger().Error("the request body could not be decoded")
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}
