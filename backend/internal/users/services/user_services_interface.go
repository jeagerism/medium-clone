package services

import "github.com/jeagerism/medium-clone/backend/internal/users/entities"

type UserService interface {
	GetUserProfile(id int) (*entities.UserWithStats, error)
	AddFollowing(req entities.UserAddFollowingRequest) error
	DeleteFollowing(req entities.UserAddFollowingRequest) error
	Login(req entities.LoginRequest) (*entities.UserProfileResponse, error)
	Register(req entities.RegisterRequest) (*entities.UserWithStats, error)
	RefreshAccessToken(refreshToken string) (*entities.UserToken, error)
}
