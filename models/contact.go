package models

import (
	"github.com/shoksin/go-contacts-REST-API-/pkg/logging"
	u "github.com/shoksin/go-contacts-REST-API-/utils"

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
	logging.GetLogger().Info("The contact has been successfully confirmed")
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
		logging.GetLogger().Error("couldn't get a contact from the database")
		return nil
	}
	return contact
}

func GetContacts(id uint) []*Contact {
	contacts := make([]*Contact, 0)
	err := GetDB().Table("contacts").Where("user_id = ?", id).Find(&contacts).Error
	if err != nil {
		logging.GetLogger().Error("couldn't get a contacts from the database")
		return nil
	}
	return contacts
}

func DeleteContact(userId uint, name string) map[string]interface{} {
	err := GetDB().Table("contacts").Where("user_id = ? AND name = ?", userId, name).Delete(&Contact{}).Error
	if err != nil {
		logging.GetLogger().Error("couldn't delete a contact from the database")
		return u.Message(false, "deletion error")
	}
	return u.Message(true, "The contact was successfully deleted")
}

func DeleteContacts(id uint) map[string]interface{} {
	err := GetDB().Table("contacts").Where("user_id = ?", id).Delete(&Contact{}).Error
	if err != nil {
		logging.GetLogger().Error("couldn't delete a contact from the database")
		return u.Message(false, "deletion error")
	}
	return u.Message(true, "All contacts were deleted")
}

func UpdateContact(user_id uint, contact_id uint, cont *Contact) map[string]interface{} {
	contact := &Contact{}
	err := GetDB().Table("contacts").Where("id = ? AND user_id = ?", contact_id, user_id).First(&contact).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logging.GetLogger().Info("The contact was not registered in the database")
			return u.Message(false, "The contact not found")
		}
		logging.GetLogger().Error("Error when trying to get a contact from the database")
		return u.Message(false, "Error with DB")
	}
	contact.Name = cont.Name
	contact.Phone = cont.Phone

	if err := GetDB().Save(&contact).Error; err != nil {
		logging.GetLogger().Error("Error when saving a contact in the database")
		return u.Message(false, "update error")
	}
	return u.Message(true, "contact updated successfully")
}

func PatchContact(user_id uint, contact_id uint, cont *Contact) map[string]interface{} {
	contact := &Contact{}
	err := GetDB().Table("contacts").Where("id = ? AND user_id = ?", contact_id, user_id).First(&contact).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Contact not found")
		}
		logging.GetLogger().Error("Error when trying to get a contact from the database")
		return u.Message(false, "update error")
	}
	if cont.Name != "" {
		contact.Name = cont.Name
	}
	if cont.Phone != "" {
		contact.Phone = cont.Phone
	}

	if err := GetDB().Save(&contact).Error; err != nil {
		logging.GetLogger().Error("Contact wasn't saved")
		return u.Message(false, "update error")
	}
	return u.Message(true, "success")
}
