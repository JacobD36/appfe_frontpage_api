package dto

import (
	"errors"
	"regexp"

	"github.com/JacobD36/appfe_frontpage_api/pkg/validator"
)

type CreateUserInput struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role,omitempty"`
}

func (c *CreateUserInput) Validate() error {
	if err := validator.Validate.Struct(c); err != nil {
	}

	pass := c.Password
	if len(pass) < 8 {
		return errors.New(ErrPasswordMinLength)
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pass)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(pass)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(pass)
	hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>\/?\\|]`).MatchString(pass)
	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New(ErrPasswordComplexity)
	}

	return nil
}
