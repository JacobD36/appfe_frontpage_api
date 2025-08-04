package dto

import "github.com/JacobD36/appfe_frontpage_api/internal/domain"

type AuthLoginResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}
