package usecase

import (
	"context"
	"fmt"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"github.com/JacobD36/appfe_frontpage_api/pkg/logger"
)

// messagingService implementa el servicio interno de mensajería
type messagingService struct {
	emailService interfaces.EmailService
	logger       logger.Logger
}

// NewMessagingService crea una nueva instancia del servicio de mensajería
func NewMessagingService(emailService interfaces.EmailService, logger logger.Logger) interfaces.MessagingService {
	return &messagingService{
		emailService: emailService,
		logger:       logger,
	}
}

// SendEmail envía un correo electrónico simple
func (s *messagingService) SendEmail(ctx context.Context, to, subject, htmlContent string) error {
	// Validar parámetros
	if to == "" {
		return fmt.Errorf(dto.ErrMessagingRecipientRequired)
	}
	if subject == "" {
		return fmt.Errorf(dto.ErrMessagingSubjectRequired)
	}
	if htmlContent == "" {
		return fmt.Errorf(dto.ErrMessagingContentRequired)
	}

	// Validar formato de email usando el servicio
	if err := s.emailService.ValidateEmail(to); err != nil {
		return fmt.Errorf(dto.ErrMessagingInvalidEmailFormat, err)
	}

	// Log del intento de envío
	s.logger.Info(ctx, dto.MsgMessagingSendingEmail,
		logger.String("to", to),
		logger.String("subject", subject),
	)

	// Enviar email usando el proveedor
	emailMessage := &interfaces.EmailMessage{
		To:          []string{to},
		Subject:     subject,
		HtmlContent: htmlContent,
		TextContent: "", // Por defecto vacío, solo HTML
	}

	err := s.emailService.SendEmail(ctx, emailMessage)
	if err != nil {
		s.logger.Error(ctx, dto.MsgMessagingFailedToSendEmail,
			logger.Error("error", err),
			logger.String("to", to),
			logger.String("subject", subject),
		)
		return fmt.Errorf(dto.ErrMessagingFailedToSend, err)
	}

	// Log de éxito
	s.logger.Info(ctx, dto.MsgMessagingEmailSentSuccess,
		logger.String("to", to),
		logger.String("subject", subject),
	)

	return nil
}
