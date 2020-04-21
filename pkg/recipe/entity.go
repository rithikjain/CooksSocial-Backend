package recipe

import "github.com/jinzhu/gorm"

type Recipe struct {
	gorm.Model
	UserID      uint   `json:"-"`
	RecipeName  string `json:"recipe_name"`
	Description string `json:"description"`
	Difficulty  uint   `json:"difficulty"`
	Procedure   string `json:"procedure"`
	ImgUrl      string `json:"img_url"`
	Likes       uint   `json:"likes"`
}
