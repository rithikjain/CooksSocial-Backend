package recipe

import (
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/rithikjain/SocialRecipe/pkg/entities"
)

type Service interface {
	CreateRecipe(recipe *entities.Recipe) (*entities.Recipe, error)

	UpdateRecipe(recipe *entities.Recipe) (*entities.Recipe, error)

	FindUserByID(id uint) (*entities.User, error)

	FindRecipeByID(recipeID uint) (*entities.Recipe, error)

	LikeRecipe(userID, recipeID uint) error

	UnlikeRecipe(userID, recipeID uint) error

	ShowUsersWhoLiked(recipeID uint, pageNo int) (*pagination.Paginator, error)

	GetAllRecipesOfUser(userID uint, pageNo int) (*pagination.Paginator, error)

	ShowUsersFavRecipes(userID uint, pageNo int) (*pagination.Paginator, error)

	ShowUserFeed(userID uint, pageNo int) (*pagination.Paginator, error)

	DeleteRecipe(recipeID uint) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) CreateRecipe(recipe *entities.Recipe) (*entities.Recipe, error) {
	return s.repo.CreateRecipe(recipe)
}

func (s *service) UpdateRecipe(recipe *entities.Recipe) (*entities.Recipe, error) {
	return s.repo.UpdateRecipe(recipe)
}

func (s *service) FindUserByID(id uint) (*entities.User, error) {
	return s.repo.FindUserByID(id)
}

func (s *service) FindRecipeByID(recipeID uint) (*entities.Recipe, error) {
	return s.repo.FindRecipeByID(recipeID)
}

func (s *service) LikeRecipe(userID, recipeID uint) error {
	return s.repo.LikeRecipe(userID, recipeID)
}

func (s *service) UnlikeRecipe(userID, recipeID uint) error {
	return s.repo.UnlikeRecipe(userID, recipeID)
}

func (s *service) ShowUsersWhoLiked(recipeID uint, pageNo int) (*pagination.Paginator, error) {
	return s.repo.GetUsersWhoLiked(recipeID, pageNo)
}

func (s *service) GetAllRecipesOfUser(userID uint, pageNo int) (*pagination.Paginator, error) {
	return s.repo.GetAllRecipesOfUser(userID, pageNo)
}

func (s *service) ShowUsersFavRecipes(userID uint, pageNo int) (*pagination.Paginator, error) {
	return s.repo.GetUsersFavRecipes(userID, pageNo)
}

func (s *service) ShowUserFeed(userID uint, pageNo int) (*pagination.Paginator, error) {
	return s.repo.ShowUserFeed(userID, pageNo)
}

func (s *service) DeleteRecipe(recipeID uint) error {
	return s.repo.DeleteRecipe(recipeID)
}
