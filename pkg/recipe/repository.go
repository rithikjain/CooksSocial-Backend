package recipe

import (
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/jinzhu/gorm"
	"github.com/rithikjain/SocialRecipe/pkg"
	"github.com/rithikjain/SocialRecipe/pkg/entities"
)

type Repository interface {
	CreateRecipe(recipe *entities.Recipe) (*entities.Recipe, error)

	UpdateRecipe(recipe *entities.Recipe) (*entities.Recipe, error)

	FindUserByID(id uint) (*entities.User, error)

	FindRecipeByID(recipeID uint) (*entities.Recipe, error)

	LikeRecipe(userID, recipeID uint) error

	UnlikeRecipe(userID, recipeID uint) error

	GetUsersWhoLiked(recipeID uint, pageNo int) (*pagination.Paginator, error)

	GetAllRecipesOfUser(userID uint, pageNo int) (*pagination.Paginator, error)

	GetUsersFavRecipes(userID uint, pageNo int) (*pagination.Paginator, error)

	ShowUserFeed(userID uint, pageNo int) (*pagination.Paginator, error)

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

func (r *repo) FindUserByID(id uint) (*entities.User, error) {
	user := &entities.User{}
	r.DB.Where("id = ?", id).First(user)
	if user.Email == "" {
		return nil, pkg.ErrNotFound
	}
	return user, nil
}

func (r *repo) FindRecipeByID(recipeID uint) (*entities.Recipe, error) {
	recipe := &entities.Recipe{}
	err := r.DB.Where("id = ?", recipeID).First(recipe).Error
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	return recipe, nil
}

func (r *repo) LikeRecipe(userID, recipeID uint) error {
	recipe := &entities.Recipe{}
	err := r.DB.Where("id = ?", recipeID).First(recipe).Error
	if err != nil {
		return pkg.ErrDatabase
	}
	recipe.Likes += 1
	result := r.DB.Save(recipe)
	if result.Error != nil {
		return pkg.ErrDatabase
	}

	like := &entities.LikeDetail{
		RecipeID: recipeID,
		UserID:   userID,
	}
	err = r.DB.Create(like).Error
	if err != nil {
		return pkg.ErrDatabase
	}

	return nil
}

func (r *repo) UnlikeRecipe(userID, recipeID uint) error {
	recipe := &entities.Recipe{}
	err := r.DB.Where("id = ?", recipeID).First(recipe).Error
	if err != nil {
		return pkg.ErrDatabase
	}
	recipe.Likes -= 1
	result := r.DB.Save(recipe)
	if result.Error != nil {
		return pkg.ErrDatabase
	}
	err = r.DB.Where("user_id = ? and recipe_id = ?", userID, recipeID).Unscoped().Delete(&entities.LikeDetail{}).Error
	if err != nil {
		return pkg.ErrDatabase
	}

	return nil
}

func (r *repo) GetUsersWhoLiked(recipeID uint, pageNo int) (*pagination.Paginator, error) {
	var likes []entities.LikeDetail
	if err := r.DB.Where("recipe_id = ?", recipeID).Find(&likes).Error; err != nil {
		return nil, pkg.ErrDatabase
	}

	var userIDs []uint
	for _, like := range likes {
		userIDs = append(userIDs, like.UserID)
	}

	var users []entities.User
	stmt := r.DB.Where(userIDs)
	page := pagination.Paging(&pagination.Param{
		DB:      stmt,
		Page:    pageNo,
		Limit:   10,
		OrderBy: []string{"created_at desc"},
	}, &users)
	return page, nil
}

func (r *repo) GetAllRecipesOfUser(userID uint, pageNo int) (*pagination.Paginator, error) {
	var recipes []entities.Recipe
	stmt := r.DB.Where("user_id = ?", userID)
	page := pagination.Paging(&pagination.Param{
		DB:      stmt,
		Page:    pageNo,
		Limit:   7,
		OrderBy: []string{"created_at desc"},
	}, &recipes)
	return page, nil
}

func (r *repo) GetUsersFavRecipes(userID uint, pageNo int) (*pagination.Paginator, error) {
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
	stmt := r.DB.Where(favouriteRecipeIDs)
	page := pagination.Paging(&pagination.Param{
		DB:      stmt,
		Page:    pageNo,
		Limit:   7,
		OrderBy: []string{"created_at desc"},
	}, &recipes)
	return page, nil
}

func (r *repo) ShowUserFeed(userID uint, pageNo int) (*pagination.Paginator, error) {
	var followings []entities.Following
	if err := r.DB.Where("user_id = ?", userID).Find(&followings).Error; err != nil {
		return nil, pkg.ErrDatabase
	}

	var otherUserIDs []uint
	for _, following := range followings {
		otherUserIDs = append(otherUserIDs, following.OthersUserID)
	}

	var recipes []entities.Recipe
	stmt := r.DB.Where("user_id in (?)", otherUserIDs)
	page := pagination.Paging(&pagination.Param{
		DB:      stmt,
		Page:    pageNo,
		Limit:   7,
		OrderBy: []string{"created_at desc"},
	}, &recipes)
	return page, nil
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
