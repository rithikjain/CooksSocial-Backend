package user

import (
	"github.com/jinzhu/gorm"
	"github.com/rithikjain/SocialRecipe/pkg"
	"github.com/rithikjain/SocialRecipe/pkg/entities"
)

type Repository interface {
	FindByID(id uint) (*entities.User, error)

	FindByEmail(email string) (*entities.User, error)

	Register(user *entities.User) (*entities.User, error)

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

func (r *repo) FindByID(id uint) (*entities.User, error) {
	user := &entities.User{}
	r.DB.Where("id = ?", id).First(user)
	if user.Email == "" {
		return nil, pkg.ErrNotFound
	}
	return user, nil
}

func (r *repo) Register(user *entities.User) (*entities.User, error) {
	result := r.DB.Save(user)
	if result.Error != nil {
		return nil, pkg.ErrDatabase
	}
	return user, nil
}

func (r *repo) DoesEmailExist(email string) (bool, error) {
	user := &entities.User{}
	if r.DB.Where("email = ?", email).First(user).RecordNotFound() {
		return false, nil
	}
	return true, nil
}

func (r *repo) FindByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	result := r.DB.Where("email = ?", email).First(user)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, pkg.ErrNotFound
	}
	return user, nil
}

func (r *repo) AddRecipeToFav(userID, recipeID uint) error {
	fav := &entities.FavoriteRecipe{
		UserID:   userID,
		RecipeID: recipeID,
	}
	if err := r.DB.Create(fav).Error; err != nil {
		return pkg.ErrDatabase
	}
	return nil
}

func (r *repo) RemoveRecipeFromFav(userID, recipeID uint) error {
	err := r.DB.Where("user_id = ? and recipe_id = ?", userID, recipeID).Unscoped().Delete(&entities.FavoriteRecipe{}).Error
	if err != nil {
		return pkg.ErrDatabase
	}
	return nil
}
