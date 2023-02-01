package auth

import (
	"context"
	"errors"
	"fmt"

	"ghostorange/internal/pkg/argon2hash"
)



type (
	UsrStorage interface{
		GetPasswordHash(cxt context.Context, userName string) (string, string, error)
		AddUser(ctx context.Context, username, hash string) (string, error)
		UserExists(ctx context.Context, username string) (bool, error)
	
	}
)

func Login(ctx context.Context, userName, password string, p UsrStorage) (string, error) {
	userID, pwdHash, err := p.GetPasswordHash(ctx, userName)
	if err != nil {
		return "", err
	} else if userID == "" {
		return "", ErrUnathorized
	}

	ok, err := argon2hash.ComparePasswordAndHash(password, pwdHash)
	if err != nil {
		return "", fmt.Errorf("failed to verify password: %w", err)
	}

	if !ok {
		return "", ErrUnathorized
	}

	return userID, nil
}

func RegisterUser(ctx context.Context, userName, password string, us UsrStorage) (string, error) {
	err := validateUserName(ctx, userName, us)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			return "", fmt.Errorf("invalid user name %w", err)
		}

		return "", err
	}

	hash, err := argon2hash.GenerateFromPassword(password, argon2hash.DefaultParams())

	if err != nil {
		return "", fmt.Errorf("failed to generate hash from password: %w", err)
	}

	return us.AddUser(ctx, userName, hash)
}

func validateUserName(ctx context.Context, userName string, us UsrStorage) error {
	if exists, err := us.UserExists(ctx, userName); err != nil {
		return fmt.Errorf("failed to check if user already exists %w", err)
	} else if exists {
		return ErrUserAlreadyExists
	}

	return nil
}
