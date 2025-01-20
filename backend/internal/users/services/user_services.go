package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jeagerism/medium-clone/backend/config"
	"github.com/jeagerism/medium-clone/backend/internal/users/entities"
	"github.com/jeagerism/medium-clone/backend/internal/users/repositories"
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

func (s *userService) Login(req entities.LoginRequest) (*entities.UserWithStats, string, error) {
	// ตรวจสอบว่าผู้ใช้งานมีอยู่ในระบบหรือไม่
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, "", ErrUserNotFound
	}

	// ตรวจสอบความถูกต้องของรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", ErrInvalidPassword
	}

	// เตรียมข้อมูล Claims สำหรับ JWT
	claims := jwt.MapClaims{
		"email":   req.Email,
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(s.config.GetAccessTokenExpiry()).Unix(), // ใช้ค่า timeout จาก Config
	}

	// สร้าง JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := s.config.GetJWTSecret() // ดึง Secret Key จาก Config
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, "", ErrGenToken
	}

	// ดึงข้อมูลผู้ใช้งานเพิ่มเติม
	userProfile, err := s.userRepo.FindUser(user.ID)
	if err != nil {
		return nil, "", ErrUserNotFound
	}

	// ส่งคืนข้อมูลโปรไฟล์ผู้ใช้พร้อม JWT Token
	return userProfile, t, nil
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
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user", // ค่าเริ่มต้นสำหรับ Role
		Bio:      req.Bio,
	}

	// บันทึกผู้ใช้งานใหม่ลงในฐานข้อมูล
	userID, err := s.userRepo.CreateUser(newUser)
	if err != nil {
		return nil, ErrCreateUserFailed
	}

	// ดึงข้อมูลโปรไฟล์ผู้ใช้งานใหม่
	createdUser, err := s.userRepo.FindUser(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return createdUser, nil
}
