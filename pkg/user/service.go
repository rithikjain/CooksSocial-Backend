package user

import (
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/rithikjain/SocialRecipe/pkg"
	"github.com/rithikjain/SocialRecipe/pkg/entities"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type Service interface {
	Register(user *entities.User) (*entities.User, error)

	Login(email, password string) (*entities.User, error)

	DoesEmailExist(email string) (bool, error)

	DoesUsernameExist(username string) (bool, error)

	GetUserByID(id uint) (*entities.User, error)

	AddRecipeToFav(userID, recipeID uint) error

	RemoveRecipeFromFav(userID, recipeID uint) error

	FollowUser(userID, otherUserID uint) error

	UnFollowUser(userID, otherUserID uint) error

	ViewFollowers(userID uint, pageNo int) (*pagination.Paginator, error)

	ViewFollowing(userID uint, pageNo int) (*pagination.Paginator, error)

	SearchUsers(userID uint, query string, pageNo int) (*pagination.Paginator, error)

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

func Validate(user *entities.User) (bool, error) {
	if !strings.Contains(user.Email, "@") {
		return false, pkg.ErrEmail
	}

	if len(user.Password) < 6 || len(user.Password) > 60 {
		return false, pkg.ErrPassword
	}
	return true, nil
}

func (s *service) DoesEmailExist(email string) (bool, error) {
	return s.repo.DoesEmailExist(email)
}

func (s *service) DoesUsernameExist(username string) (bool, error) {
	return s.repo.DoesUsernameExist(username)
}

func (s *service) Register(user *entities.User) (*entities.User, error) {
	// Validation
	validate, err := Validate(user)
	if !validate {
		return nil, err
	}

	pass, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = pass
	return s.repo.Register(user)
}

func (s *service) Login(email, password string) (*entities.User, error) {
	user := &entities.User{}
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if CheckPasswordHash(password, user.Password) {
		return user, nil
	}
	return nil, pkg.ErrNotFound
}

func (s *service) GetUserByID(id uint) (*entities.User, error) {
	return s.repo.FindByID(id)
}

func (s *service) AddRecipeToFav(userID, recipeID uint) error {
	return s.repo.AddRecipeToFav(userID, recipeID)
}

func (s *service) RemoveRecipeFromFav(userID, recipeID uint) error {
	return s.repo.RemoveRecipeFromFav(userID, recipeID)
}

func (s *service) FollowUser(userID, otherUserID uint) error {
	return s.repo.FollowUser(userID, otherUserID)
}

func (s *service) UnFollowUser(userID, otherUserID uint) error {
	return s.repo.UnFollowUser(userID, otherUserID)
}

func (s *service) ViewFollowers(userID uint, pageNo int) (*pagination.Paginator, error) {
	return s.repo.ViewFollowers(userID, pageNo)
}

func (s *service) ViewFollowing(userID uint, pageNo int) (*pagination.Paginator, error) {
	return s.repo.ViewFollowing(userID, pageNo)
}

func (s *service) SearchUsers(userID uint, query string, pageNo int) (*pagination.Paginator, error) {
	return s.repo.SearchUsers(userID, query, pageNo)
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
