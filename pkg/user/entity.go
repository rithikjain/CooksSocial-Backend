package user

import (
	"github.com/jinzhu/gorm"
	"github.com/rithikjain/SocialRecipe/pkg/recipe"
)

type User struct {
	gorm.Model
	Name          string          `json:"name"`
	PhoneNumber   string          `json:"phone_number"`
	Email         string          `json:"email"`
	Password      string          `json:"password"`
	ProfileImgUrl string          `json:"profile_img_url"`
	Bio           string          `json:"bio"`
	Recipes       []recipe.Recipe `json:"-" gorm:"foreignkey:UserID"`
}
