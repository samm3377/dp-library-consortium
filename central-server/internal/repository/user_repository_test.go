package repository

import (
	"central-server/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	assert.NoError(t, err)

	err = db.AutoMigrate(&models.User{})

	assert.NoError(t, err)

	return db
}

func TestUserRepository_Create_Success(t *testing.T) {

	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Username:  "test",
		Email:     "test@test.com",
		Password:  "hashed_password",
		CreatedAt: time.Now(),
	}

	err := repo.Create(user)

	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	var saved models.User

	err = db.First(&saved, user.ID).Error

	assert.NoError(t, err)
	assert.Equal(t, "test", saved.Username)
	assert.Equal(t, "test@test.com", saved.Email)
}

func TestUserRepository_FindByUsername_Success(t *testing.T) {

	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: "password",
	}

	db.Create(user)
	result, err := repo.FindByUsername("test")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test", result.Username)
	assert.Equal(t, "test@test.com", result.Email)
}

func TestUserRepository_FindByUsername_NotFound(t *testing.T) {

	db := setupTestDB(t)
	repo := NewUserRepository(db)
	user, err := repo.FindByUsername("unknown")

	assert.Error(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "", user.Username)
}
