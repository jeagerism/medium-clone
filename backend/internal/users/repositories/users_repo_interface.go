package repositories

import "github.com/jeagerism/medium-clone/backend/internal/users/entities"

type UserRepository interface {
	FindUser(id int) (*entities.UserWithStats, error)
	SaveFollowing(req entities.UserAddFollowingRequest) error
	RemoveFollowing(req entities.UserAddFollowingRequest) error
	GetUserByEmail(email string) (*entities.UserCredentials, error)
	CreateUser(user entities.User) (int, error)
}
