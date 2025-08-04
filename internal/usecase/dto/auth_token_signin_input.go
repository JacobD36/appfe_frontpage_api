package dto

type AuthTokenSignInInput struct {
	Token string `json:"token" validate:"required" example:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
