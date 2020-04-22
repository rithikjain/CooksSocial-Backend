package recipe

type Service interface {
	CreateRecipe(recipe *Recipe) (*Recipe, error)

	UpdateRecipe(recipe *Recipe) (*Recipe, error)

	FindRecipeByID(recipeID uint) (*Recipe, error)

	LikeRecipe(recipeID uint) (*Recipe, error)

	UnlikeRecipe(recipeID uint) (*Recipe, error)

	GetAllRecipesOfUser(userID uint) (*[]Recipe, error)

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

func (s *service) CreateRecipe(recipe *Recipe) (*Recipe, error) {
	return s.repo.CreateRecipe(recipe)
}

func (s *service) UpdateRecipe(recipe *Recipe) (*Recipe, error) {
	return s.repo.UpdateRecipe(recipe)
}

func (s *service) FindRecipeByID(recipeID uint) (*Recipe, error) {
	return s.repo.FindRecipeByID(recipeID)
}

func (s *service) LikeRecipe(recipeID uint) (*Recipe, error) {
	return s.repo.LikeRecipe(recipeID)
}

func (s *service) UnlikeRecipe(recipeID uint) (*Recipe, error) {
	return s.repo.UnlikeRecipe(recipeID)
}

func (s *service) GetAllRecipesOfUser(userID uint) (*[]Recipe, error) {
	return s.repo.GetAllRecipesOfUser(userID)
}

func (s *service) DeleteRecipe(recipeID uint) error {
	return s.repo.DeleteRecipe(recipeID)
}
