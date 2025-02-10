package repositories

import (
	"time"

	"github.com/jeagerism/medium-clone/backend/internal/users/entities"
)

type UserRepository interface {
	FindUser(id int) (*entities.UserWithStats, error)
	SaveFollowing(req entities.UserAddFollowingRequest) error
	RemoveFollowing(req entities.UserAddFollowingRequest) error
	GetUserByEmail(email string) (*entities.UserCredentials, error)
	CreateUser(user entities.User) (int, error)
	SaveRefreshToken(userID int, token string, expiresAt time.Time) error
	UpdateRefreshToken(userID int, newToken string, expiresAt time.Time) error
	GetRefresh(token string) (*entities.UserCredentials, error)
	DeleteRefreshToken(userID int) error
}
