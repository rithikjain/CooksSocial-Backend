package recipe

import "github.com/rithikjain/SocialRecipe/pkg/entities"

type Service interface {
	CreateRecipe(recipe *entities.Recipe) (*entities.Recipe, error)

	UpdateRecipe(recipe *entities.Recipe) (*entities.Recipe, error)

	FindRecipeByID(recipeID uint) (*entities.Recipe, error)

	LikeRecipe(recipeID uint) (*entities.Recipe, error)

	UnlikeRecipe(recipeID uint) (*entities.Recipe, error)

	GetAllRecipesOfUser(userID uint) (*[]entities.Recipe, error)

	ShowUsersFavRecipes(userID uint) (*[]entities.Recipe, error)

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

func (s *service) FindRecipeByID(recipeID uint) (*entities.Recipe, error) {
	return s.repo.FindRecipeByID(recipeID)
}

func (s *service) LikeRecipe(recipeID uint) (*entities.Recipe, error) {
	return s.repo.LikeRecipe(recipeID)
}

func (s *service) UnlikeRecipe(recipeID uint) (*entities.Recipe, error) {
	return s.repo.UnlikeRecipe(recipeID)
}

func (s *service) GetAllRecipesOfUser(userID uint) (*[]entities.Recipe, error) {
	return s.repo.GetAllRecipesOfUser(userID)
}

func (s *service) ShowUsersFavRecipes(userID uint) (*[]entities.Recipe, error) {
	return s.repo.GetUsersFavRecipes(userID)
}

func (s *service) DeleteRecipe(recipeID uint) error {
	return s.repo.DeleteRecipe(recipeID)
}
