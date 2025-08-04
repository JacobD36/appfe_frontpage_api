package domain

import (
	"errors"
	"time"
)

const (
	UserRole       = "USER_ROLE"
	AdminRole      = "ADMIN_ROLE"
	ErrInvalidRole = "rol inválido. Los roles válidos son: USER_ROLE, ADMIN_ROLE"
)

var ValidRoles = []string{UserRole, AdminRole}

func IsValidRole(role string) bool {
	for _, validRole := range ValidRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

func ValidateRole(role string) (string, error) {
	if role == "" {
		return UserRole, nil
	}

	if !IsValidRole(role) {
		return "", errors.New(ErrInvalidRole)
	}

	return role, nil
}

type User struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Password       *string    `json:"password"`
	Img            *string    `json:"img,omitempty"`
	Role           string     `json:"role"`
	Status         bool       `json:"status"`
	EmailValidated bool       `json:"emailValidated"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
