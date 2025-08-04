package middleware

import (
	"net/http"
	"strings"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain"
	"github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"github.com/labstack/echo/v4"
)

type JWTMiddleware struct {
	jwtService interfaces.JWTService
}

func NewJWTMiddleware(jwtService interfaces.JWTService) *JWTMiddleware {
	return &JWTMiddleware{
		jwtService: jwtService,
	}
}

func (m *JWTMiddleware) Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"code":    http.StatusUnauthorized,
					"message": dto.ErrTokenMissing,
					"status":  "Unauthorized",
				})
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"code":    http.StatusUnauthorized,
					"message": dto.ErrInvalidTokenFormat,
					"status":  "Unauthorized",
				})
			}

			tokenString := tokenParts[1]

			claims, err := m.jwtService.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"code":    http.StatusUnauthorized,
					"message": dto.ErrInvalidToken,
					"status":  "Unauthorized",
				})
			}

			c.Set("user_id", claims.ID)
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}

func (m *JWTMiddleware) RequiredRole(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("user_role").(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"code":    http.StatusUnauthorized,
					"message": dto.ErrTokenMissing,
					"status":  "Unauthorized",
				})
			}

			if userRole != requiredRole {
				return c.JSON(http.StatusForbidden, map[string]any{
					"code":    http.StatusForbidden,
					"message": dto.ErrInsufficientPermissions,
					"status":  "Forbidden",
				})
			}

			return next(c)
		}
	}
}

func (m *JWTMiddleware) RequireUserRole() echo.MiddlewareFunc {
	return m.RequiredRole(domain.UserRole)
}

func (m *JWTMiddleware) RequireAdminRole() echo.MiddlewareFunc {
	return m.RequiredRole(domain.AdminRole)
}

func (m *JWTMiddleware) RequireAnyRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("user_role").(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"code":    http.StatusUnauthorized,
					"message": dto.ErrTokenMissing,
					"status":  "Unauthorized",
				})
			}

			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]any{
				"code":    http.StatusForbidden,
				"message": dto.ErrInsufficientPermissions,
				"status":  "Forbidden",
			})
		}
	}
}
