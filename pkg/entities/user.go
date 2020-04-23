package entities

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name               string           `json:"name"`
	Username           string           `json:"username"`
	PhoneNumber        string           `json:"phone_number"`
	Email              string           `json:"email"`
	Password           string           `json:"password"`
	ProfileImgUrl      string           `json:"profile_img_url"`
	ProfileImgPublicID string           `json:"-"`
	Bio                string           `json:"bio"`
	Recipes            []Recipe         `json:"-" gorm:"foreignkey:UserID"`
	FavouriteRecipes   []FavoriteRecipe `json:"-" gorm:"foreignkey:UserID"`
	Following          []Following      `json:"-" gorm:"foreignkey:UserID"`
	Followers          []Follower       `json:"-" gorm:"foreignkey:UserID"`
}

type FavoriteRecipe struct {
	gorm.Model
	UserID   uint
	RecipeID uint
}

type Following struct {
	gorm.Model
	UserID       uint
	OthersUserID uint
}

type Follower struct {
	gorm.Model
	UserID       uint
	OthersUserID uint
}
