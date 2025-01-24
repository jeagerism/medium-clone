package services

import (
	"fmt"
	"time"

	"github.com/jeagerism/medium-clone/backend/config"
	"github.com/jeagerism/medium-clone/backend/internal/users/entities"
	"github.com/jeagerism/medium-clone/backend/internal/users/repositories"
	"github.com/jeagerism/medium-clone/backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repositories.UserRepository
	config   *config.Config
}

func NewUserService(userRepo repositories.UserRepository, cfg *config.Config) UserService {
	return &userService{
		userRepo: userRepo,
		config:   cfg, // เก็บ config ไว้ใน struct ของ service
	}
}

func (s *userService) GetUserProfile(id int) (*entities.UserWithStats, error) {
	user, err := s.userRepo.FindUser(id)
	fmt.Println(id)
	if err != nil {
		return &entities.UserWithStats{}, err
	}
	return user, nil
}

func (s *userService) AddFollowing(req entities.UserAddFollowingRequest) error {
	if err := s.userRepo.SaveFollowing(req); err != nil {
		return ErrFailedToFollow
	}
	return nil
}

func (s *userService) DeleteFollowing(req entities.UserAddFollowingRequest) error {
	if err := s.userRepo.RemoveFollowing(req); err != nil {
		return ErrFailedToDeleteFollowing
	}
	return nil
}

func (s *userService) Login(req entities.LoginRequest) (*entities.UserProfileResponse, error) {
	// ตรวจสอบว่าผู้ใช้งานมีอยู่ในระบบหรือไม่
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// ตรวจสอบความถูกต้องของรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	// สร้าง Access Token
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Role,
		string(s.config.JWT().GetJWTSecret()),
		s.config.JWT().GetAccessTokenExpiry(),
	)
	if err != nil {
		return nil, ErrGenToken
	}

	// สร้าง Refresh Token
	refreshToken, err := utils.GenerateRefreshToken(
		user.ID,
		string(s.config.JWT().GetJWTSecret()),
		s.config.JWT().GetRefreshTokenExpiry(),
	)
	if err != nil {
		return nil, ErrGenToken
	}

	// บันทึก Refresh Token ในฐานข้อมูล
	err = s.userRepo.SaveRefreshToken(user.ID, refreshToken, time.Now().Add(s.config.JWT().GetRefreshTokenExpiry()))
	if err != nil {
		return nil, ErrGenToken
	}

	// ส่งคืนข้อมูล User พร้อม Access Token และ Refresh Token
	return &entities.UserProfileResponse{
		User: &entities.User{
			ID:    user.ID,
			Email: req.Email,
			Role:  user.Role,
		},
		Token: &entities.UserToken{
			ID:           fmt.Sprintf("%d", user.ID),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *userService) Register(req entities.RegisterRequest) (*entities.UserWithStats, error) {
	// ตรวจสอบว่าอีเมลนี้มีอยู่ในระบบหรือไม่
	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// เข้ารหัสรหัสผ่าน
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrHashPassword
	}

	// สร้างข้อมูลผู้ใช้งานใหม่
	newUser := entities.User{
		Name:         req.Name,
		Email:        req.Email,
		Password:     string(hashedPassword),
		Role:         "user", // ค่าเริ่มต้นสำหรับ Role
		Bio:          req.Bio,
		ProfileImage: req.ProfileImage,
	}

	// บันทึกผู้ใช้งานใหม่ลงในฐานข้อมูล
	userID, err := s.userRepo.CreateUser(newUser)
	if err != nil {
		return nil, ErrCreateUserFailed
	}

	return s.GetUserProfile(userID)
}

func (s *userService) RefreshAccessToken(refreshToken string) (*entities.UserToken, error) {
	// ตรวจสอบ Refresh Token
	_, claims, err := utils.ValidateToken(refreshToken, string(s.config.JWT().GetJWTSecret()))
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// ตรวจใน DB
	user, err := s.userRepo.GetRefresh(refreshToken)
	if err != nil {
		return nil, ErrGenToken
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("user_id is missing or invalid in token claims")
	}
	userID := int(userIDFloat)

	if userID != user.ID {
		return nil, ErrUserNotFound
	}

	// สร้าง Access Token ใหม่
	newAccessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Role,
		string(s.config.JWT().GetJWTSecret()),
		s.config.JWT().GetAccessTokenExpiry(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	// สร้าง Refresh Token ใหม่
	newRefreshToken, err := utils.GenerateRefreshToken(
		user.ID,
		string(s.config.JWT().GetJWTSecret()),
		s.config.JWT().GetRefreshTokenExpiry(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	// อัปเดต Refresh Token ในฐานข้อมูล
	expiresAt := time.Now().Add(s.config.JWT().GetRefreshTokenExpiry())
	if err := s.userRepo.UpdateRefreshToken(user.ID, newRefreshToken, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	// ส่งคืน Token ใหม่
	return &entities.UserToken{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
