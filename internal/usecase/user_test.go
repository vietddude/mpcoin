package usecase

import (
	"context"
	"errors"
	"mpc/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock type for the UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, params domain.CreateHashedUserParams) (domain.User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUser(ctx context.Context, id int64) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUseCase(mockRepo)

	ctx := context.Background()
	createParams := domain.CreateUserParams{
		Email:    "test@example.com",
		Password: "password123",
	}

	t.Run("Successful user creation", func(t *testing.T) {
		expectedUser := domain.User{
			ID:    1,
			Email: "test@example.com",
		}

		mockRepo.On("GetUserByEmail", ctx, createParams.Email).Return(domain.User{}, errors.New("user not found"))
		mockRepo.On("CreateUser", ctx, mock.AnythingOfType("domain.CreateHashedUserParams")).Return(expectedUser, nil)

		user, err := uc.CreateUser(ctx, createParams)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User already exists", func(t *testing.T) {
		existingUser := domain.User{
			ID:    1,
			Email: "test@example.com",
		}

		mockRepo.On("GetUserByEmail", ctx, createParams.Email).Return(existingUser, nil)

		_, err := uc.CreateUser(ctx, createParams)

		assert.Error(t, err)
		assert.EqualError(t, err, "user with this email already exists")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error creating user", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", ctx, createParams.Email).Return(domain.User{}, errors.New("user not found"))
		mockRepo.On("CreateUser", ctx, mock.AnythingOfType("domain.CreateHashedUserParams")).Return(domain.User{}, errors.New("database error"))

		_, err := uc.CreateUser(ctx, createParams)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertExpectations(t)
	})
}

// Add more test functions as needed
