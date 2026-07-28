package service

import (
	"central-server/internal/models"
	"central-server/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(username, email, password string, time time.Time) error
	Login(username, password string) (*models.User, error)
}

type authService struct {
	userDB repository.UserRepository
}

func NewAuthService(userDB repository.UserRepository) AuthService {
	return &authService{userDB}
}

func (s *authService) Register(username, email, password string, time time.Time) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{Username: username, Email: email, Password: string(hashedPassword), CreatedAt: time}
	return s.userDB.Create(user)
}

func (s *authService) Login(username, password string) (*models.User, error) {
	user, err := s.userDB.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
