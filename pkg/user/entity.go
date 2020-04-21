package user

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Bio         string `json:"bio"`
}
