package user

import (
	"pkg/database"
	"pkg/errors"
)

type Repository interface {
	Create(u *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	FindAllByRole(role string, out *[]*User) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(u *User) error {
	return database.DB.Create(u).Error
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var u User
	if err := database.DB.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, errors.Wrap(err, "find user by email")
	}
	return &u, nil
}

func (r *repository) FindByID(id string) (*User, error) {
	var u User
	if err := database.DB.First(&u, "id = ?", id).Error; err != nil {
		return nil, errors.Wrap(err, "find user by id")
	}
	return &u, nil
}

func (r *repository) FindAllByRole(role string, out *[]*User) error {
	return database.DB.Where("role = ?", role).Find(out).Error
}
