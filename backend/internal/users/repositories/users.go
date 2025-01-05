package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/jeagerism/medium-clone/backend/internal/users/entities"
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
