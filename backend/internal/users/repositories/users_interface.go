package repositories

import "github.com/jeagerism/medium-clone/backend/internal/users/entities"

type UserRepository interface {
	FindUser(id int) (*entities.UserWithStats, error)
}
