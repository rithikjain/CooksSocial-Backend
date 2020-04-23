package user

import (
	"github.com/rithikjain/SocialRecipe/pkg"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type Service interface {
	Register(user *User) (*User, error)

	Login(email, password string) (*User, error)

	GetUserByID(id uint) (*User, error)

	AddRecipeToFav(userID, recipeID uint) error

	RemoveRecipeFromFav(userID, recipeID uint) error

	GetRepo() Repository
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (user *User) Validate() (bool, error) {
	if !strings.Contains(user.Email, "@") {
		return false, pkg.ErrEmail
	}

	if len(user.Password) < 6 || len(user.Password) > 60 {
		return false, pkg.ErrPassword
	}
	return true, nil
}

func (s *service) Register(user *User) (*User, error) {
	// Validation
	validate, err := user.Validate()
	if !validate {
		return nil, err
	}

	exists, err := s.repo.DoesEmailExist(user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		//noinspection GoErrorStringFormat
		return nil, pkg.ErrExists
	}
	pass, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = pass
	return s.repo.Register(user)
}

func (s *service) Login(email, password string) (*User, error) {
	user := &User{}
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if CheckPasswordHash(password, user.Password) {
		return user, nil
	}
	return nil, pkg.ErrNotFound
}

func (s *service) GetUserByID(id uint) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *service) AddRecipeToFav(userID, recipeID uint) error {
	return s.repo.AddRecipeToFav(userID, recipeID)
}

func (s *service) RemoveRecipeFromFav(userID, recipeID uint) error {
	return s.repo.RemoveRecipeFromFav(userID, recipeID)
}

func (s *service) GetRepo() Repository {
	return s.repo
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
