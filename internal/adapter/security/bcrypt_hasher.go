package security

import (
	"errors"
	"fmt"

	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

var ErrEmptyPassword = errors.New(dto.ErrEmptyPassword)

func NewBcryptHasher(cost int) *BcryptHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}

	if len(password) > 72 {
		return "", errors.New(dto.ErrPasswordTooLong)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf(dto.ErrHashGeneration, err)
	}

	return string(hash), nil
}

func (h *BcryptHasher) Verify(hashedPassword, password string) error {
	if hashedPassword == "" {
		return errors.New(dto.ErrHashEmpty)
	}

	if password == "" {
		return ErrEmptyPassword
	}

	if len(password) > 72 {
		return errors.New(dto.ErrPasswordTooLong)
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New(dto.ErrPasswordIncorrect)
		}
		return fmt.Errorf(dto.ErrPasswordVerification, err)
	}

	return nil
}

func (h *BcryptHasher) GetCost() int {
	return h.cost
}

func (h *BcryptHasher) IsValidHash(hash string) bool {
	if len(hash) < 60 {
		return false
	}

	if len(hash) < 4 || hash[0] != '$' {
		return false
	}

	validPrefixes := []string{"$2a$", "$2b$", "$2x$", "$2y$"}
	for _, prefix := range validPrefixes {
		if len(hash) >= len(prefix) && hash[:len(prefix)] == prefix {
			return true
		}
	}

	return false
}
