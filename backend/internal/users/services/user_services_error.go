package services

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrFailedToFollow = errors.New("can not follow")
)
