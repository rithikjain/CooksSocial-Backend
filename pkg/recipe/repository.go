package recipe

type Repository interface {
	SaveRecipe(recipe *Recipe) (*Recipe, error)

	FindRecipeByID(recipeID uint) (*Recipe, error)

	DeleteRecipe(recipeID uint) error
}

//func ()
