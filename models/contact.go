package models

import (
	u "go-contacts/utils"

	"github.com/jinzhu/gorm"
)

type Contact struct {
	gorm.Model
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	UserID uint   `json:"user_id"`
}

func (contact *Contact) Validate() (map[string]interface{}, bool) {
	if contact.Name == "" {
		return u.Message(false, "Invalid username!"), false
	}
	if contact.Phone == "" {
		return u.Message(false, "Invalid user phone!"), false
	}
	if contact.UserID < 0 {
		return u.Message(false, "Invalid userID"), false
	}

	//все обязательные поля присутствуют
	return u.Message(true, "success"), true
}

func (contact *Contact) Create() map[string]interface{} {
	if resp, ok := contact.Validate(); !ok {
		return resp
	}

	GetDB().Create(contact)

	resp := u.Message(true, "success")
	resp["contact"] = contact
	return resp
}

func GetContact(id uint) *Contact {
	contact := &Contact{}
	err := GetDB().Table("contacts").Where("id = ?", id).First(contact).Error
	if err != nil {
		return nil
	}
	return contact
}

func GetContacts(id uint) []*Contact {
	contacts := make([]*Contact, 0)
	//err := GetDB().Table("contacts").Where("user_id = ?", id).Find(&contacts).Error
	err := GetDB().Table("contacts").Where("user_id = ?", id).Find(&contacts).Error
	if err != nil {
		return nil
	}
	return contacts
}
