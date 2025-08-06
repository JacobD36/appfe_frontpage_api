package usecase

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain"
	ui "github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/interfaces"
)

type userService struct {
	uowFactory       ui.UnitOfWorkFactory
	hasher           ui.PasswordHasher
	messagingService ui.MessagingService
	templateService  ui.TemplateService
}

func NewUserService(
	uowFactory ui.UnitOfWorkFactory,
	h ui.PasswordHasher,
	messagingService ui.MessagingService,
	templateService ui.TemplateService) interfaces.UserService {
	return &userService{
		uowFactory:       uowFactory,
		hasher:           h,
		messagingService: messagingService,
		templateService:  templateService,
	}
}

func (s *userService) Create(ctx context.Context, u *domain.User) error {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.Rollback()

	validatedRole, err := domain.ValidateRole(u.Role)
	if err != nil {
		return err
	}
	u.Role = validatedRole

	u.CreatedAt = time.Now()
	u.Status = true
	u.Name = strings.ToUpper(strings.TrimSpace(u.Name))
	u.EmailValidated = true

	// Capturar la contraseña original antes del hash para enviarla por email
	var originalPassword string
	if u.Password != nil {
		originalPassword = *u.Password
		hashed, err := s.hasher.Hash(*u.Password)
		if err != nil {
			return err
		}
		u.Password = &hashed
	}

	if err := uow.UserRepository().Create(ctx, u); err != nil {
		return err
	}

	// Enviar email de bienvenida después de crear el usuario
	if s.messagingService != nil && s.templateService != nil {
		welcomeContent, err := s.templateService.RenderWelcomeEmail(u.Name, originalPassword)
		if err != nil {
			// Log del error pero no fallar la creación del usuario
			return uow.Commit()
		}

		// Enviar email de forma asíncrona para no bloquear la respuesta
		go func() {
			if err := s.messagingService.SendEmail(context.Background(), u.Email, dto.WelcomeEmailSubject, welcomeContent); err != nil {
				// Log del error pero no fallar la creación del usuario
				// El logger debería estar disponible en el servicio de mensajería
			}
		}()
	}

	return uow.Commit()
}

func (s *userService) UpdateByID(ctx context.Context, input ui.UpdateUserInput) error {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.Rollback()

	if err := uow.UserRepository().UpdateByID(ctx, input); err != nil {
		return err
	}
	return uow.Commit()
}

func (s *userService) GetAll(ctx context.Context, pagination *domain.Pagination) (*domain.PaginatedResult[*domain.User], error) {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return nil, err
	}
	defer uow.Rollback()

	users, total, err := uow.UserRepository().GetAll(ctx, pagination)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		user.Password = nil
	}

	if pagination == nil {
		pagination = domain.NewPagination(1, int(total), "")
	}

	result := domain.NewPaginatedResult(users, pagination, total)
	return result, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return nil, err
	}
	defer uow.Rollback()

	user, err := uow.UserRepository().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return nil, err
	}
	defer uow.Rollback()

	return uow.UserRepository().FindByEmail(ctx, email)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.Rollback()

	if err := uow.UserRepository().Delete(ctx, id); err != nil {
		return err
	}
	return uow.Commit()
}

func (s *userService) CreateInitialAdmin(ctx context.Context) error {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminName := os.Getenv("ADMIN_NAME")

	existingAdmin, err := s.FindByEmail(ctx, adminEmail)
	if err != nil && err.Error() != dto.ErrNoRowsFound {
		return err
	}

	if existingAdmin != nil {
		return nil
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		return errors.New("ADMIN_PASSWORD environment variable is required")
	}

	hashedPassword, err := s.hasher.Hash(adminPassword)
	if err != nil {
		return err
	}

	adminUser := &domain.User{
		Name:           adminName,
		Email:          adminEmail,
		Password:       &hashedPassword,
		Role:           domain.AdminRole,
		Status:         true,
		EmailValidated: true,
		CreatedAt:      time.Now(),
	}

	uow, err := s.uowFactory.New(ctx)
	if err != nil {
		return err
	}
	defer uow.Rollback()

	if err := uow.UserRepository().Create(ctx, adminUser); err != nil {
		return err
	}

	return uow.Commit()
}
