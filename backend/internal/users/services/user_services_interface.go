package services

import "github.com/jeagerism/medium-clone/backend/internal/users/entities"

type UserService interface {
	GetUserProfile(id int) (*entities.UserWithStats, error)
}
