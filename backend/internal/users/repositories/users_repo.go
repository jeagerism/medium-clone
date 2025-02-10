package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jeagerism/medium-clone/backend/internal/users/entities"
	"github.com/jeagerism/medium-clone/backend/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindUser(id int) (*entities.UserWithStats, error) {
	query := `
    SELECT 
        u.id,
        u."name",
        u.bio,
        u.profile_image,
        COUNT(DISTINCT f.follower_id) AS follower_count,
        jsonb_agg(
            jsonb_build_object('id', fu.id, 'name', fu.name)
        ) AS following
    FROM users u
    LEFT JOIN follows f ON u.id = f.follower_id
    LEFT JOIN users fu ON f.following_id = fu.id
    WHERE u.id = $1
    GROUP BY u.id, u."name", u.bio, u.profile_image;
	`

	// Define a raw struct that maps to the query result.
	type UserRaw struct {
		ID            int    `db:"id"`
		Name          string `db:"name"`
		Bio           string `db:"bio"`
		ProfileImage  string `db:"profile_image"`
		FollowerCount int    `db:"follower_count"`
		Following     string `db:"following"` // Change to string to handle raw JSON as text
	}

	// Initialize a variable to hold the raw query result
	var raw UserRaw
	err := r.db.Get(&raw, query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying user data: %w", err)
	}

	// Map the raw result into the UserWithStats struct
	var user entities.UserWithStats
	user.ID = raw.ID
	user.Name = raw.Name
	user.Bio = raw.Bio
	user.ProfileImage = raw.ProfileImage
	user.FollowerCount = raw.FollowerCount

	// Unmarshal the `following` field from JSON into `Follow` structs
	if len(raw.Following) > 0 {
		var follows []entities.Follow
		// Convert the raw `following` JSON string into a valid JSON array
		if err := json.Unmarshal([]byte(raw.Following), &follows); err != nil {
			return nil, fmt.Errorf("error unmarshalling following data: %w", err)
		}
		user.Following = follows
	}

	// Return the user with stats
	return &user, nil
}

func (r *userRepository) SaveFollowing(req entities.UserAddFollowingRequest) error {
	query := `INSERT INTO follows (follower_id,following_id) VALUES ($1,$2);`
	_, err := r.db.Exec(query, req.FollowerID, req.FollowingID)
	if err != nil {
		logger.LogError(fmt.Errorf("failed to save following '%d' & '%d': %w", req.FollowerID, req.FollowingID, err))
		return fmt.Errorf("failed to save following '%d' & '%d': %w", req.FollowerID, req.FollowingID, err)
	}
	return nil
}

func (r *userRepository) RemoveFollowing(req entities.UserAddFollowingRequest) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2;`
	_, err := r.db.Exec(query, req.FollowerID, req.FollowingID)
	if err != nil {
		logger.LogError(fmt.Errorf("failed to remove following '%d' & '%d': %w", req.FollowerID, req.FollowingID, err))
		return fmt.Errorf("failed to remove following '%d' & '%d': %w", req.FollowerID, req.FollowingID, err)
	}
	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*entities.UserCredentials, error) {
	logger.LogInfo(fmt.Sprintf("Querying user by email: %s", email))

	var userCreds entities.UserCredentials
	query := `SELECT id, email, password_hash, role FROM users WHERE LOWER(email) = LOWER($1);`

	err := r.db.Get(&userCreds, query, email)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed to retrieve user credentials for email %s: %w", email, err))
		return nil, fmt.Errorf("failed to retrieve user credentials: %w", err)
	}

	logger.LogInfo(fmt.Sprintf("Successfully retrieved user credentials for email: %s", email))
	return &userCreds, nil
}

func (r *userRepository) CreateUser(user entities.User) (int, error) {
	logger.LogInfo(fmt.Sprintf("Creating user: %+v", user))

	query := `
		INSERT INTO users (name, email, password_hash, role, bio,profile_image)
		VALUES ($1, $2, $3, $4, $5,$6)
		RETURNING id;
	`

	var userID int
	err := r.db.QueryRow(query, user.Name, user.Email, user.Password, user.Role, user.Bio, user.ProfileImage).Scan(&userID)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed to create user %+v: %w", user, err))
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	logger.LogInfo(fmt.Sprintf("Successfully created user with ID: %d", userID))
	return userID, nil
}

func (r *userRepository) SaveRefreshToken(userID int, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO tokens (user_id, refresh_token, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`

	_, err := r.db.Exec(query, userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

func (r *userRepository) UpdateRefreshToken(userID int, newToken string, expiresAt time.Time) error {
	query := `
		UPDATE tokens
		SET refresh_token = $1, expires_at = $2, updated_at = NOW()
		WHERE user_id = $3;
	`

	_, err := r.db.Exec(query, newToken, expiresAt, userID)
	if err != nil {
		return fmt.Errorf("failed to update refresh token for user ID %d: %w", userID, err)
	}

	return nil
}

func (r *userRepository) GetRefresh(token string) (*entities.UserCredentials, error) {
	var user entities.UserCredentials
	query := `
	SELECT
		u.id ,
		u.email ,
		u."role"
	FROM users u 
	LEFT JOIN tokens t on u.id = t.user_id 
	WHERE t.refresh_token = $1;
	`
	err := r.db.Get(&user, query, token)
	if err != nil {
		logger.LogError(fmt.Errorf("Failed to retrieve user credentials for refresh token %s: %w", token, err))
		return nil, fmt.Errorf("failed to retrieve user credentials: %w", err)
	}
	return &user, nil
}

func (r *userRepository) DeleteRefreshToken(userID int) error {
	query := `DELETE FROM tokens WHERE user_id = $1;`
	_, err := r.db.Exec(query, userID)
	return err
}
