package entities

import "github.com/jinzhu/gorm"

type Recipe struct {
	gorm.Model
	UserID      uint         `json:"user_id"`
	RecipeName  string       `json:"recipe_name"`
	Description string       `json:"description"`
	Ingredients string       `json:"ingredients"`
	Difficulty  int          `json:"difficulty"`
	Procedure   string       `json:"procedure"`
	ImgUrl      string       `json:"img_url"`
	ImgPublicId string       `json:"-"`
	Username    string       `json:"username"`
	UserImg     string       `json:"user_img"`
	Likes       int          `json:"likes"`
	LikeDetails []LikeDetail `json:"-" gorm:"foreignkey:RecipeID"`
}

type LikeDetail struct {
	gorm.Model
	RecipeID uint
	UserID   uint
}
