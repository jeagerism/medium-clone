package services

import (
	"fmt"

	"github.com/jeagerism/medium-clone/backend/internal/users/entities"
	"github.com/jeagerism/medium-clone/backend/internal/users/repositories"
)

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserProfile(id int) (*entities.UserWithStats, error) {
	user, err := s.userRepo.FindUser(id)
	fmt.Println(id)
	if err != nil {
		return &entities.UserWithStats{}, err
	}
	return user, nil
}
