package user

import (
	"github.com/jinzhu/gorm"
	"github.com/rithikjain/SocialRecipe/pkg"
)

type Repository interface {
	FindByID(id uint) (*User, error)

	FindByEmail(email string) (*User, error)

	Register(user *User) (*User, error)

	DoesEmailExist(email string) (bool, error)

	AddRecipeToFav(userID, recipeID uint) error

	RemoveRecipeFromFav(userID, recipeID uint) error
}

type repo struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repo{
		DB: db,
	}
}

func (r *repo) FindByID(id uint) (*User, error) {
	user := &User{}
	r.DB.Where("id = ?", id).First(user)
	if user.Email == "" {
		return nil, pkg.ErrNotFound
	}
	return user, nil
}

func (r *repo) Register(user *User) (*User, error) {
	result := r.DB.Save(user)
	if result.Error != nil {
		return nil, pkg.ErrDatabase
	}
	return user, nil
}

func (r *repo) DoesEmailExist(email string) (bool, error) {
	user := &User{}
	if r.DB.Where("email = ?", email).First(user).RecordNotFound() {
		return false, nil
	}
	return true, nil
}

func (r *repo) FindByEmail(email string) (*User, error) {
	user := &User{}
	result := r.DB.Where("email = ?", email).First(user)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, pkg.ErrNotFound
	}
	return user, nil
}

func (r *repo) AddRecipeToFav(userID, recipeID uint) error {
	user := &User{}
	r.DB.Where("id = ?", userID).First(user)
	if user.Email == "" {
		return pkg.ErrNotFound
	}
	for _, id := range user.FavouriteRecipeIDs {
		if id == recipeID {
			return pkg.ErrExists
		}
	}
	user.FavouriteRecipeIDs = append(user.FavouriteRecipeIDs, recipeID)
	return nil
}

func (r *repo) RemoveRecipeFromFav(userID, recipeID uint) error {
	user := &User{}
	r.DB.Where("id = ?", userID).First(user)
	if user.Email == "" {
		return pkg.ErrNotFound
	}
	for idx, id := range user.FavouriteRecipeIDs {
		if id == recipeID {
			user.FavouriteRecipeIDs = RemoveIndex(user.FavouriteRecipeIDs, idx)
			return nil
		}
	}
	return pkg.ErrNotFound
}

func RemoveIndex(s []uint, index int) []uint {
	return append(s[:index], s[index+1:]...)
}
