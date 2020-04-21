package user

import (
	"github.com/jinzhu/gorm"
	"github.com/rithikjain/SocialRecipe/pkg"
)

type Repository interface {
	FindByID(id float64) (*User, error)

	FindByEmail(email string) (*User, error)

	Register(user *User) (*User, error)

	DoesEmailExist(email string) (bool, error)
}

type repo struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repo{
		DB: db,
	}
}

func (r *repo) FindByID(id float64) (*User, error) {
	user := &User{}
	r.DB.Where("id = ?", id).First(user)
	if user.Email == "" {
		return nil, pkg.ErrNotFound
	}
	return user, nil
}

func (r *repo) Register(user *User) (*User, error) {
	result := r.DB.Save(user)
	if result.Error != nil {
		return nil, pkg.ErrDatabase
	}
	return user, nil
}

func (r *repo) DoesEmailExist(email string) (bool, error) {
	user := &User{}
	if r.DB.Where("email = ?", email).First(user).RecordNotFound() {
		return false, nil
	}
	return true, nil
}

func (r *repo) FindByEmail(email string) (*User, error) {
	user := &User{}
	result := r.DB.Where("email = ?", email).First(user)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, pkg.ErrNotFound
	}
	return user, nil
}
