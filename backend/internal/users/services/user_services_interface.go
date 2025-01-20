package services

import "github.com/jeagerism/medium-clone/backend/internal/users/entities"

type UserService interface {
	GetUserProfile(id int) (*entities.UserWithStats, error)
	AddFollowing(req entities.UserAddFollowingRequest) error
	DeleteFollowing(req entities.UserAddFollowingRequest) error
	Login(req entities.LoginRequest) (*entities.UserWithStats, string, error)
	Register(req entities.RegisterRequest) (*entities.UserWithStats, error)
}
