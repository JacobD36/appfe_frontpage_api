package handler

import (
	"net/http"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/pkg/logger"
	"github.com/JacobD36/appfe_frontpage_api/pkg/validator"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService interfaces.UserService
}

func NewUserHandler(userService interfaces.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var input dto.CreateUserInput

	logger.Info(ctx, dto.MsgCreatingNewUser, logger.String("operation", "user_create"))

	if err := c.Bind(&input); err != nil {
		logger.Warn(ctx, dto.MsgInvalidInputForUser, logger.Error("error", err))
		return Error(c, http.StatusBadRequest, dto.ErrInvalidInput)
	}

	if err := input.Validate(); err != nil {
		msg := dto.TranslateValidationErrors(err)
		logger.Warn(ctx, dto.MsgUserValidationFailed,
			logger.String("email", input.Email),
			logger.Error("validation_error", err),
		)
		return Error(c, http.StatusBadRequest, msg)
	}

	existing, err := h.userService.FindByEmail(ctx, input.Email)
	if err != nil && err.Error() != dto.ErrNoRowsFound {
		logger.LogError(ctx, dto.MsgErrorCheckingExistingUser,
			logger.String("email", input.Email),
			logger.Error("error", err),
		)
		return Error(c, http.StatusInternalServerError, dto.ErrInternalServer)
	}

	if existing != nil {
		logger.Warn(ctx, dto.MsgAttemptCreateExistingUser,
			logger.String("email", input.Email),
		)
		return Error(c, http.StatusConflict, dto.ErrUserAlreadyExists)
	}

	u := &domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: &input.Password,
		Role:     input.Role,
	}

	if err := h.userService.Create(ctx, u); err != nil {
		logger.LogError(ctx, dto.MsgFailedToCreateUser,
			logger.String("email", input.Email),
			logger.String("role", input.Role),
			logger.Error("error", err),
		)
		return Error(c, http.StatusInternalServerError, err.Error())
	}

	logger.Info(ctx, dto.MsgUserCreatedSuccessfully,
		logger.String("user_id", u.ID),
		logger.String("email", u.Email),
		logger.String("role", u.Role),
	)

	return Success(c, http.StatusCreated, dto.ErrUserCreatedSuccess, echo.Map{
		dto.UserIdLabel:    u.ID,
		dto.UserEmailLabel: u.Email,
	})
}

func (h *UserHandler) GetAll(c echo.Context) error {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")
	searchStr := c.QueryParam("search")

	var pagination *domain.Pagination = nil
	var err error

	if pageStr != "" || limitStr != "" || searchStr != "" {
		pagination, err = domain.ParsePaginationFromQuery(pageStr, limitStr, searchStr)
		if err != nil {
			return Error(c, http.StatusBadRequest, err.Error())
		}
	}

	ctx := c.Request().Context()

	result, err := h.userService.GetAll(ctx, pagination)
	if err != nil {
		return Error(c, http.StatusInternalServerError, dto.ErrInternalServer)
	}

	message := dto.ErrUserRetrievedSuccess
	if pagination != nil {
		if pagination.Search != "" {
			message = dto.ErrUsersSearchSuccess
		} else if pagination.Page > 1 || pagination.Limit < 100 {
			message = dto.ErrUsersRetrievedPaginated
		}
	}

	return Success(c, http.StatusOK, message, result)
}

func (h *UserHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return Error(c, http.StatusBadRequest, dto.ErrInvalidUserID)
	}

	ctx := c.Request().Context()

	user, err := h.userService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == dto.ErrNoRowsFound {
			return Error(c, http.StatusNotFound, dto.ErrUserNotFound)
		}
		return Error(c, http.StatusInternalServerError, dto.ErrInternalServer)
	}

	user.Password = nil
	return Success(c, http.StatusOK, dto.ErrUserRetrievedSuccess, user)
}

func (h *UserHandler) UpdateByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return Error(c, http.StatusBadRequest, dto.ErrInvalidUserID)
	}

	var input dto.UpdateUserInput
	if err := c.Bind(&input); err != nil {
		return Error(c, http.StatusBadRequest, dto.ErrInvalidInput)
	}

	if err := validator.Validate.Struct(input); err != nil {
		return Error(c, http.StatusBadRequest, dto.TranslateValidationErrors(err))
	}

	input.ID = id

	ctx := c.Request().Context()

	_, err := h.userService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == dto.ErrNoRowsFound {
			return Error(c, http.StatusNotFound, dto.ErrUserNotFound)
		}
		return Error(c, http.StatusInternalServerError, dto.ErrInternalServer)
	}

	if err := h.userService.UpdateByID(ctx, &input); err != nil {
		return Error(c, http.StatusInternalServerError, err.Error())
	}

	return Success(c, http.StatusOK, dto.ErrUserUpdatedSuccess, echo.Map{
		dto.UserIdLabel: id,
	})
}

func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return Error(c, http.StatusBadRequest, dto.ErrInvalidUserID)
	}

	ctx := c.Request().Context()

	_, err := h.userService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == dto.ErrNoRowsFound {
			return Error(c, http.StatusNotFound, dto.ErrUserNotFound)
		}
		return Error(c, http.StatusInternalServerError, dto.ErrInternalServer)
	}

	if err := h.userService.Delete(ctx, id); err != nil {
		return Error(c, http.StatusInternalServerError, err.Error())
	}

	return Success(c, http.StatusOK, dto.ErrUserDeletedSuccess, echo.Map{
		dto.UserIdLabel: id,
	})
}
