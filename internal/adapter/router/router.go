package router

import (
	"github.com/JacobD36/appfe_frontpage_api/internal/adapter/handler"
	"github.com/JacobD36/appfe_frontpage_api/internal/adapter/middleware"
	domainInterfaces "github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	usecaseInterfaces "github.com/JacobD36/appfe_frontpage_api/internal/usecase/interfaces"
	v "github.com/JacobD36/appfe_frontpage_api/pkg/validator"
	validatorLib "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Router struct {
	e        *echo.Echo
	handlers *Handlers
	jwtMw    *middleware.JWTMiddleware
}

type Handlers struct {
	User usecaseInterfaces.UserService
	Auth usecaseInterfaces.AuthService
}

type CustomValidator struct {
	validator *validatorLib.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func New(userService usecaseInterfaces.UserService, authService usecaseInterfaces.AuthService, jwtService domainInterfaces.JWTService) *Router {
	e := echo.New()

	e.Use(
		echoMiddleware.Recover(),
		middleware.Logger(), // Usar nuestro logger personalizado
		echoMiddleware.Secure(),
		echoMiddleware.Gzip(),
		echoMiddleware.CORS(),
	)

	e.Validator = &CustomValidator{validator: v.Validate}

	jwtMw := middleware.NewJWTMiddleware(jwtService)

	router := &Router{
		e:     e,
		jwtMw: jwtMw,
		handlers: &Handlers{
			User: userService,
			Auth: authService,
		},
	}

	router.registerRoutes()

	return router
}

func (r *Router) registerRoutes() {
	v1 := r.e.Group("/api/v1")

	userHandler := handler.NewUserHandler(r.handlers.User)
	userGroup := v1.Group("/users")

	adminUserGroup := userGroup.Group("", r.jwtMw.Authenticate(), r.jwtMw.RequireAdminRole())
	adminUserGroup.POST("", userHandler.Create)
	adminUserGroup.GET("", userHandler.GetAll)
	adminUserGroup.GET("/:id", userHandler.GetByID)
	adminUserGroup.PUT("/:id", userHandler.UpdateByID)
	adminUserGroup.DELETE("/:id", userHandler.Delete)

	authHandler := handler.NewAuthHandler(r.handlers.Auth)
	authGroup := v1.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/sign-in-with-token", authHandler.SignInWithToken)
}

func (r *Router) Start(addr string) error {
	return r.e.Start(addr)
}

func (r *Router) GetEchoInstance() *echo.Echo {
	return r.e
}
