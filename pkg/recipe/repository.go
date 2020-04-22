package recipe

import (
	"github.com/jinzhu/gorm"
	"github.com/rithikjain/SocialRecipe/pkg"
)

type Repository interface {
	CreateRecipe(recipe *Recipe) (*Recipe, error)

	UpdateRecipe(recipe *Recipe) (*Recipe, error)

	FindRecipeByID(recipeID uint) (*Recipe, error)

	LikeRecipe(recipeID uint) (*Recipe, error)

	UnlikeRecipe(recipeID uint) (*Recipe, error)

	GetAllRecipesOfUser(userID uint) (*[]Recipe, error)

	DeleteRecipe(recipeID uint) error
}

type repo struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repo{
		DB: db,
	}
}

func (r *repo) CreateRecipe(recipe *Recipe) (*Recipe, error) {
	result := r.DB.Create(recipe)
	if result.Error != nil {
		return nil, pkg.ErrDatabase
	}
	return recipe, nil
}

func (r *repo) UpdateRecipe(recipe *Recipe) (*Recipe, error) {
	result := r.DB.Save(recipe)
	if result.Error != nil {
		return nil, pkg.ErrDatabase
	}
	return recipe, nil
}

func (r *repo) FindRecipeByID(recipeID uint) (*Recipe, error) {
	recipe := &Recipe{}
	err := r.DB.Where("id = ?", recipeID).First(recipe).Error
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	return recipe, nil
}

func (r *repo) LikeRecipe(recipeID uint) (*Recipe, error) {
	recipe := &Recipe{}
	err := r.DB.Where("id = ?", recipeID).First(recipe).Error
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	recipe.Likes += 1
	result := r.DB.Save(recipe)
	if result.Error != nil {
		return nil, pkg.ErrDatabase
	}
	return recipe, nil
}

func (r *repo) UnlikeRecipe(recipeID uint) (*Recipe, error) {
	recipe := &Recipe{}
	err := r.DB.Where("id = ?", recipeID).First(recipe).Error
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	recipe.Likes += 1
	result := r.DB.Save(recipe)
	if result.Error != nil {
		return nil, pkg.ErrDatabase
	}
	return recipe, nil
}

func (r *repo) GetAllRecipesOfUser(userID uint) (*[]Recipe, error) {
	var recipes []Recipe
	err := r.DB.Where("user_id = ?", userID).Find(&recipes).Error
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	return &recipes, nil
}

func (r *repo) DeleteRecipe(recipeID uint) error {
	recipe := &Recipe{}
	err := r.DB.Where("id = ?", recipeID).First(recipe).Error
	if err != nil {
		return pkg.ErrDatabase
	}
	err = r.DB.Delete(recipe).Error
	if err != nil {
		return pkg.ErrDatabase
	}
	return nil
}
