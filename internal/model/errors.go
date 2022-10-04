package model

import "errors"

var (
	ErrUserNotFound            = errors.New("user doesn't exists")
	ErrVerificationCodeInvalid = errors.New("verification code is invalid")
	ErrProjectNotFound         = errors.New("project doesn't exists")
	ErrUserIsAlreadyExists     = errors.New("user with this email already exists")
	ErrUnknownCallbackType     = errors.New("unknown callback type")
	ErrProjectIsAlreadyExists  = errors.New("this project is already exists")
	ErrCouldParseId            = errors.New("ID of users is not in correct format")
)
