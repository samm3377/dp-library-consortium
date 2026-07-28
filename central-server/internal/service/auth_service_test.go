package service

import (
	"central-server/internal/models"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.User), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	now := time.Now()

	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	err := service.Register(
		"test",
		"test@test.com",
		"password123",
		now,
	)

	assert.NoError(t, err)

	mockRepo.AssertCalled(t, "Create", mock.MatchedBy(func(u *models.User) bool {
		if u.Username != "test" {
			return false
		}
		if u.Email != "test@test.com" {
			return false
		}
		if u.CreatedAt != now {
			return false
		}
		return bcrypt.CompareHashAndPassword(
			[]byte(u.Password),
			[]byte("password123"),
		) == nil
	}))
}

func TestRegister_CreateFails(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)
	expectedErr := errors.New("db error")
	mockRepo.On("Create", mock.Anything).Return(expectedErr)

	err := service.Register(
		"test",
		"test@test.com",
		"password123",
		time.Now(),
	)

	assert.Equal(t, expectedErr, err)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	user := &models.User{
		Username: "test",
		Password: string(hashed),
	}

	mockRepo.On("FindByUsername", "test").Return(user, nil)
	result, err := service.Login("test", "password123")

	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)
	expectedErr := errors.New("user not found")
	mockRepo.On("FindByUsername", "test").Return(nil, expectedErr)
	user, err := service.Login("test", "password123")

	assert.Nil(t, user)
	assert.Equal(t, expectedErr, err)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)
	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctPassword"), bcrypt.DefaultCost)

	user := &models.User{
		Username: "test",
		Password: string(hashed),
	}

	mockRepo.On("FindByUsername", "test").Return(user, nil)
	result, err := service.Login("test", "wrongPassword")

	assert.Nil(t, result)
	assert.Error(t, err)
}
