package repositories_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jeagerism/medium-clone/backend/internal/users/entities"
	"github.com/jeagerism/medium-clone/backend/internal/users/repositories"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestFindUser(t *testing.T) {
	// Set up sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	// Wrap the mock *sql.DB with *sqlx.DB
	sqlxDB := sqlx.NewDb(db, "postgres")

	// Create the user repository instance with the mocked database
	repo := repositories.NewUserRepository(sqlxDB)

	// Define expected user data
	userID := 3
	expectedUser := &entities.UserWithStats{
		ID:            userID,
		Name:          "John Doe",
		Bio:           "Software Developer",
		ProfileImage:  "https://example.com/profile.jpg",
		FollowerCount: 150,
		Following: []entities.Follow{
			{ID: 2, Name: "Jane Smith"},
			{ID: 5, Name: "Michael Johnson"},
		},
	}

	// Define the test cases
	tests := []struct {
		name        string
		mockSetup   func()
		expected    *entities.UserWithStats
		expectedErr string
	}{
		{
			name: "Success",
			mockSetup: func() {
				// Set up the mock query and expected result for a successful query
				rows := sqlmock.NewRows([]string{"id", "name", "bio", "profile_image", "follower_count", "following"}).
					AddRow(
						userID,
						"John Doe",
						"Software Developer",
						"https://example.com/profile.jpg",
						150,
						`[{"id": 2, "name": "Jane Smith"}, {"id": 5, "name": "Michael Johnson"}]`, // Mocked JSON array for following
					)
				// Expect the query to be executed with the specific ID and return the rows
				mock.ExpectQuery(`SELECT .* FROM users u .*`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expected:    expectedUser,
			expectedErr: "",
		},
		{
			name: "DatabaseError",
			mockSetup: func() {
				// Simulate a database query error
				mock.ExpectQuery(`SELECT .* FROM users u .*`).
					WillReturnError(fmt.Errorf("database error"))
			},
			expected:    nil,
			expectedErr: "error querying user data",
		},
		{
			name: "JSONUnmarshalError",
			mockSetup: func() {
				// Simulate invalid JSON in the `following` column
				rows := sqlmock.NewRows([]string{"id", "name", "bio", "profile_image", "follower_count", "following"}).
					AddRow(
						3,
						"John Doe",
						"Software Developer",
						"https://example.com/profile.jpg",
						150,
						`invalid json`, // Invalid JSON for following
					)
				// Expect the query to be executed with the specific ID and return the rows
				mock.ExpectQuery(`SELECT .* FROM users u .*`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expected:    nil,
			expectedErr: "error unmarshalling following data",
		},
	}

	// Run all the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the mock for the current test case
			tt.mockSetup()

			// Call the FindUser method
			user, err := repo.FindUser(userID)

			// Assert that the error matches the expected error (if any)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				// Assert no error occurred for successful cases
				assert.NoError(t, err)
			}

			// Assert that the returned user matches the expected user (if no error)
			if tt.expected != nil {
				assert.NotNil(t, user)
				assert.Equal(t, tt.expected.ID, user.ID)
				assert.Equal(t, tt.expected.Name, user.Name)
				assert.Equal(t, tt.expected.Bio, user.Bio)
				assert.Equal(t, tt.expected.ProfileImage, user.ProfileImage)
				assert.Equal(t, tt.expected.FollowerCount, user.FollowerCount)
				assert.Len(t, user.Following, len(tt.expected.Following))
				assert.Equal(t, tt.expected.Following, user.Following)
			}

			// Ensure that all expected queries were executed
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unmet expectations: %v", err)
			}
		})
	}
}
