package user

import (
	"pkg/logger"
	"pkg/utils"

	"go.uber.org/zap"
)

type Service interface {
	Register(username, email, password, role string) (*User, error)
	Authenticate(email, password string) (*User, error)
	FindByID(id string) (*User, error)
	FindByRole(role string) ([]*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(username, email, password, role string) (*User, error) {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	u := &User{
		Username: username,
		Email:    email,
		Password: hash,
		Role:     role,
	}
	if err := s.repo.Create(u); err != nil {
		logger.Log.Error("failed to save user",
			zap.String("username", username),
			zap.Error(err),
		)
		return nil, err
	}
	logger.Log.Info("user created successfully",
		zap.String("email", email))
	return u, nil
}

func (s *service) Authenticate(email, password string) (*User, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := utils.CheckPassword(u.Password, password); err != nil {
		return nil, err
	}
	userDetails := *u
	userDetails.Password = "hidden"
	return &userDetails, nil
}

func (s *service) FindByID(id string) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *service) FindByRole(role string) ([]*User, error) {
	var users []*User
	if err := s.repo.FindAllByRole(role, &users); err != nil {
		return nil, err
	}
	return users, nil
}
