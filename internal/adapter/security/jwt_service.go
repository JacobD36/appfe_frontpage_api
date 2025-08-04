package security

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWTService() (*JWTService, error) {
	if signKey == nil || verifyKey == nil {
		return nil, errors.New("las claves JWT no han sido cargadas")
	}

	return &JWTService{
		privateKey: signKey,
		publicKey:  verifyKey,
	}, nil
}

func (j *JWTService) GenerateToken(user domain.User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "geekway-api",
			Subject:   user.ID,
			ID:        user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTService) ValidateToken(tokenString string) (*domain.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("método de firma inválido")
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("claims inválidos")
	}

	user := &domain.User{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	}

	return user, nil
}
