package services

import "errors"

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrFailedToFollow          = errors.New("can not follow")
	ErrFailedToDeleteFollowing = errors.New("can not delete following")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrGenToken                = errors.New("could not generate token")
	ErrEmailAlreadyExists      = errors.New("email is already registered")
	ErrHashPassword            = errors.New("failed to hash password")
	ErrCreateUserFailed        = errors.New("failed to create user")
)
