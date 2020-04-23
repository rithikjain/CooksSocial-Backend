package recipe

import (
	"github.com/jinzhu/gorm"
	"github.com/rithikjain/SocialRecipe/pkg"
	"github.com/rithikjain/SocialRecipe/pkg/entities"
)

type Repository interface {
	CreateRecipe(recipe *entities.Recipe) (*entities.Recipe, error)

	UpdateRecipe(recipe *entities.Recipe) (*entities.Recipe, error)

	FindRecipeByID(recipeID uint) (*entities.Recipe, error)

	LikeRecipe(recipeID uint) (*entities.Recipe, error)

	UnlikeRecipe(recipeID uint) (*entities.Recipe, error)

	GetAllRecipesOfUser(userID uint) (*[]entities.Recipe, error)

	GetUsersFavRecipes(userID uint) (*[]entities.Recipe, error)

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

func (r *repo) CreateRecipe(recipe *entities.Recipe) (*entities.Recipe, error) {
	result := r.DB.Create(recipe)
	if result.Error != nil {
		return nil, pkg.ErrDatabase
	}
	return recipe, nil
}

func (r *repo) UpdateRecipe(recipe *entities.Recipe) (*entities.Recipe, error) {
	result := r.DB.Save(recipe)
	if result.Error != nil {
		return nil, pkg.ErrDatabase
	}
	return recipe, nil
}

func (r *repo) FindRecipeByID(recipeID uint) (*entities.Recipe, error) {
	recipe := &entities.Recipe{}
	err := r.DB.Where("id = ?", recipeID).First(recipe).Error
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	return recipe, nil
}

func (r *repo) LikeRecipe(recipeID uint) (*entities.Recipe, error) {
	recipe := &entities.Recipe{}
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

func (r *repo) UnlikeRecipe(recipeID uint) (*entities.Recipe, error) {
	recipe := &entities.Recipe{}
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

func (r *repo) GetAllRecipesOfUser(userID uint) (*[]entities.Recipe, error) {
	var recipes []entities.Recipe
	err := r.DB.Where("user_id = ?", userID).Find(&recipes).Error
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	return &recipes, nil
}

func (r *repo) GetUsersFavRecipes(userID uint) (*[]entities.Recipe, error) {
	var favs []entities.FavoriteRecipe
	err := r.DB.Where("user_id = ?", userID).Find(&favs).Error

	if err != nil {
		return nil, pkg.ErrDatabase
	}

	var favouriteRecipeIDs []uint
	for _, fav := range favs {
		favouriteRecipeIDs = append(favouriteRecipeIDs, fav.RecipeID)
	}

	var recipes []entities.Recipe
	er := r.DB.Where(favouriteRecipeIDs).Find(&recipes).Error
	if er != nil {
		return nil, pkg.ErrDatabase
	}
	return &recipes, nil
}

func (r *repo) DeleteRecipe(recipeID uint) error {
	recipe := &entities.Recipe{}
	err := r.DB.Where("id = ?", recipeID).First(recipe).Error
	if err != nil {
		return pkg.ErrDatabase
	}
	err = r.DB.Unscoped().Delete(recipe).Error
	if err != nil {
		return pkg.ErrDatabase
	}
	return nil
}
