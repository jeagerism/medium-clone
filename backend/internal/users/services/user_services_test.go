package services_test

import (
	"errors"
	"testing"

	"github.com/jeagerism/medium-clone/backend/internal/users/entities"
	"github.com/jeagerism/medium-clone/backend/internal/users/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindUser(id int) (*entities.UserWithStats, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.UserWithStats), args.Error(1)
}

// Test Service Layer
func TestGetUserProfile(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	userID := 1
	expectedUser := &entities.UserWithStats{
		ID:            userID,
		Name:          "John Doe",
		Bio:           "Software Developer",
		ProfileImage:  "https://example.com/profile.jpg",
		FollowerCount: 100,
	}

	// Define test cases
	tests := []struct {
		name        string
		mockReturn  func() // Mock return values
		expected    *entities.UserWithStats
		expectedErr error
	}{
		{
			name: "Success",
			mockReturn: func() {
				mockRepo.On("FindUser", userID).Return(expectedUser, nil)
			},
			expected:    expectedUser,
			expectedErr: nil,
		},
		{
			name: "DatabaseError",
			mockReturn: func() {
				mockRepo.On("FindUser", userID).Return(nil, errors.New("error querying user data"))
			},
			expected:    &entities.UserWithStats{}, // Expecting an empty user struct as a fallback
			expectedErr: services.ErrUserNotFound,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock behavior
			tt.mockReturn()

			// Call the service method
			user, err := service.GetUserProfile(userID)

			// Assert the results
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expected, user)

			// Assert that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}
