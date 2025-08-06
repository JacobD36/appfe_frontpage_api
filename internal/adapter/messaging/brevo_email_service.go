package messaging

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
	"github.com/JacobD36/appfe_frontpage_api/pkg/logger"
	brevo "github.com/getbrevo/brevo-go/lib"
)

// BrevoEmailService implementa EmailService usando la API de Brevo
type BrevoEmailService struct {
	client    *brevo.APIClient
	config    *interfaces.EmailServiceConfig
	logger    logger.Logger
	fromEmail string
	fromName  string
}

// NewBrevoEmailService crea una nueva instancia del servicio de email de Brevo
func NewBrevoEmailService(config *interfaces.EmailServiceConfig, log logger.Logger) (*BrevoEmailService, error) {
	if config == nil {
		return nil, errors.New(dto.ErrBrevoConfigRequired)
	}

	if strings.TrimSpace(config.APIKey) == "" {
		return nil, errors.New(dto.ErrBrevoAPIKeyRequired)
	}

	if strings.TrimSpace(config.FromEmail) == "" {
		return nil, errors.New(dto.ErrBrevoFromEmailRequired)
	}

	// Validar el formato del email
	if _, err := mail.ParseAddress(config.FromEmail); err != nil {
		return nil, fmt.Errorf(dto.ErrBrevoInvalidFromEmailFormat, err)
	}

	// Configurar cliente de Brevo
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader(dto.BrevoAPIKeyHeaderName, config.APIKey)
	client := brevo.NewAPIClient(cfg)

	service := &BrevoEmailService{
		client:    client,
		config:    config,
		logger:    log,
		fromEmail: config.FromEmail,
		fromName:  config.FromName,
	}

	// Validar la configuración con una prueba de conectividad
	if err := service.validateConfiguration(context.Background()); err != nil {
		return nil, fmt.Errorf(dto.ErrBrevoFailedValidateConfig, err)
	}

	service.logger.Info(context.Background(), dto.MsgBrevoServiceInitialized,
		logger.String(dto.LogFieldFromEmail, config.FromEmail),
		logger.String(dto.LogFieldFromName, config.FromName),
	)

	return service, nil
}

// validateConfiguration valida la configuración haciendo una llamada a la API
func (s *BrevoEmailService) validateConfiguration(ctx context.Context) error {
	// Crear un contexto con timeout para evitar esperas largas
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Intentar obtener información de la cuenta para validar la API key
	_, response, err := s.client.AccountApi.GetAccount(timeoutCtx)
	if err != nil {
		// Crear mensaje de error más descriptivo basado en el tipo de error
		var errorMsg string
		if response != nil {
			switch response.StatusCode {
			case 401:
				errorMsg = "invalid or expired API key"
			case 403:
				errorMsg = "API key does not have sufficient permissions"
			case 429:
				errorMsg = "API rate limit exceeded, please try again later"
			default:
				errorMsg = fmt.Sprintf("Brevo API returned status %d", response.StatusCode)
			}
		} else {
			// Error de conectividad
			errorMsg = "unable to connect to Brevo API (check internet connection)"
		}

		s.logger.Error(timeoutCtx, dto.MsgBrevoFailedValidateConfig,
			logger.Error(dto.LogFieldError, err),
			logger.String("details", errorMsg),
		)
		return fmt.Errorf("%s: %s", dto.ErrBrevoInvalidAPIKeyOrConnection, errorMsg)
	}

	return nil
}

// SendEmail envía un correo electrónico individual
func (s *BrevoEmailService) SendEmail(ctx context.Context, message *interfaces.EmailMessage) error {
	if message == nil {
		return errors.New(dto.ErrBrevoMessageRequired)
	}

	if err := s.validateEmailMessage(message); err != nil {
		s.logger.Error(ctx, dto.MsgBrevoInvalidEmailMessage,
			logger.Error(dto.LogFieldError, err),
			logger.String(dto.LogFieldSubject, message.Subject),
			logger.Any(dto.LogFieldTo, message.To),
		)
		return fmt.Errorf(dto.ErrBrevoInvalidEmailMessage, err)
	}

	// Construir el mensaje para Brevo
	sendSmtpEmail := s.buildBrevoEmail(message)

	// Enviar el email
	result, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, sendSmtpEmail)
	if err != nil {
		s.logger.Error(ctx, dto.MsgBrevoFailedSendEmail,
			logger.Error(dto.LogFieldError, err),
			logger.String(dto.LogFieldSubject, message.Subject),
			logger.Any(dto.LogFieldTo, message.To),
		)
		return fmt.Errorf(dto.ErrBrevoFailedSendEmail, err)
	}

	s.logger.Info(ctx, dto.MsgBrevoEmailSentSuccess,
		logger.String(dto.LogFieldMessageID, result.MessageId),
		logger.String(dto.LogFieldSubject, message.Subject),
		logger.Any(dto.LogFieldTo, message.To),
	)

	return nil
}

// SendBulkEmail envía correos electrónicos en lotes
func (s *BrevoEmailService) SendBulkEmail(ctx context.Context, messages []*interfaces.EmailMessage) error {
	if len(messages) == 0 {
		return errors.New(dto.ErrBrevoAtLeastOneMessageRequired)
	}

	if len(messages) > 50 { // Límite razonable para evitar problemas de memoria/tiempo
		return errors.New(dto.ErrBrevoTooManyMessages)
	}

	var failedMessages []string
	successCount := 0

	for i, message := range messages {
		if err := s.SendEmail(ctx, message); err != nil {
			failedMessages = append(failedMessages, fmt.Sprintf(dto.BrevoBulkMessageFormat, i, err.Error()))
		} else {
			successCount++
		}
	}

	s.logger.Info(ctx, dto.MsgBrevoBulkEmailCompleted,
		logger.Int(dto.LogFieldTotalMsgs, len(messages)),
		logger.Int(dto.LogFieldSuccessful, successCount),
		logger.Int(dto.LogFieldFailed, len(failedMessages)),
	)

	if len(failedMessages) > 0 {
		return fmt.Errorf(dto.ErrBrevoFailedSendMultipleMsg, len(failedMessages), strings.Join(failedMessages, "; "))
	}

	return nil
}

// SendTemplate envía un correo usando una plantilla
func (s *BrevoEmailService) SendTemplate(ctx context.Context, templateID int64, to []string, templateData map[string]interface{}) error {
	if templateID <= 0 {
		return errors.New(dto.ErrBrevoValidTemplateIDRequired)
	}

	if len(to) == 0 {
		return errors.New(dto.ErrBrevoAtLeastOneRecipient)
	}

	// Validar direcciones de email
	for _, email := range to {
		if _, err := mail.ParseAddress(email); err != nil {
			return fmt.Errorf(dto.ErrBrevoInvalidEmailAddress, email, err)
		}
	}

	// Construir los destinatarios
	recipients := make([]brevo.SendSmtpEmailTo, len(to))
	for i, email := range to {
		recipients[i] = brevo.SendSmtpEmailTo{
			Email: email,
		}
	}

	// Construir el mensaje con plantilla
	sendSmtpEmail := brevo.SendSmtpEmail{
		To:         recipients,
		TemplateId: templateID,
		Params:     templateData,
	}

	// Enviar el email
	result, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, sendSmtpEmail)
	if err != nil {
		s.logger.Error(ctx, dto.MsgBrevoFailedSendTemplate,
			logger.Error(dto.LogFieldError, err),
			logger.Any(dto.LogFieldTemplateID, templateID),
			logger.Any(dto.LogFieldTo, to),
		)
		return fmt.Errorf(dto.ErrBrevoFailedSendTemplate, err)
	}

	s.logger.Info(ctx, dto.MsgBrevoTemplateSentSuccess,
		logger.String(dto.LogFieldMessageID, result.MessageId),
		logger.Any(dto.LogFieldTemplateID, templateID),
		logger.Any(dto.LogFieldTo, to),
	)

	return nil
}

// ValidateEmail valida si una dirección de correo es válida
func (s *BrevoEmailService) ValidateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return errors.New(dto.ErrBrevoEmailRequired)
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf(dto.ErrBrevoInvalidEmailFormat, err)
	}

	return nil
}

// GetEmailStatus obtiene el estado de un correo enviado
func (s *BrevoEmailService) GetEmailStatus(ctx context.Context, messageID string) (*interfaces.EmailStatus, error) {
	if strings.TrimSpace(messageID) == "" {
		return nil, errors.New(dto.ErrBrevoMessageIDRequired)
	}

	// Nota: Brevo no proporciona una API directa para obtener el estado por message ID
	// En un escenario real, necesitarías implementar webhooks para rastrear el estado
	// o usar su API de eventos para consultar por rangos de fecha

	s.logger.Warn(ctx, dto.MsgBrevoGetStatusNotImplemented,
		logger.String(dto.LogFieldMessageID, messageID),
		logger.String(dto.LogFieldNote, dto.BrevoWebhookNote),
	)

	return &interfaces.EmailStatus{
		MessageID: messageID,
		Status:    dto.BrevoDefaultEmailStatus, // Estado por defecto
		SentAt:    "",
		Error:     "",
	}, nil
}

// validateEmailMessage valida un mensaje de email
func (s *BrevoEmailService) validateEmailMessage(message *interfaces.EmailMessage) error {
	if len(message.To) == 0 {
		return errors.New(dto.ErrBrevoAtLeastOneRecipient)
	}

	// Validar direcciones de email
	for _, email := range message.To {
		if _, err := mail.ParseAddress(email); err != nil {
			return fmt.Errorf(dto.ErrBrevoInvalidEmailAddress, email, err)
		}
	}

	if strings.TrimSpace(message.Subject) == "" {
		return errors.New(dto.ErrBrevoSubjectRequired)
	}

	if strings.TrimSpace(message.HtmlContent) == "" && strings.TrimSpace(message.TextContent) == "" {
		return errors.New(dto.ErrBrevoContentRequired)
	}

	return nil
}

// buildBrevoEmail construye un objeto SendSmtpEmail de Brevo
func (s *BrevoEmailService) buildBrevoEmail(message *interfaces.EmailMessage) brevo.SendSmtpEmail {
	// Construir destinatarios
	to := make([]brevo.SendSmtpEmailTo, len(message.To))
	for i, email := range message.To {
		to[i] = brevo.SendSmtpEmailTo{
			Email: email,
		}
	}

	// Usar el remitente configurado o el del mensaje
	fromEmail := s.fromEmail
	fromName := s.fromName

	if message.From.Email != "" {
		fromEmail = message.From.Email
		if message.From.Name != "" {
			fromName = message.From.Name
		}
	}

	sendSmtpEmail := brevo.SendSmtpEmail{
		To:      to,
		Subject: message.Subject,
		Sender: &brevo.SendSmtpEmailSender{
			Email: fromEmail,
			Name:  fromName,
		},
	}

	// Agregar contenido HTML si está presente
	if strings.TrimSpace(message.HtmlContent) != "" {
		sendSmtpEmail.HtmlContent = message.HtmlContent
	}

	// Agregar contenido de texto si está presente
	if strings.TrimSpace(message.TextContent) != "" {
		sendSmtpEmail.TextContent = message.TextContent
	}

	// Agregar ReplyTo si está presente
	if message.ReplyTo != nil && message.ReplyTo.Email != "" {
		sendSmtpEmail.ReplyTo = &brevo.SendSmtpEmailReplyTo{
			Email: message.ReplyTo.Email,
			Name:  message.ReplyTo.Name,
		}
	}

	// Agregar archivos adjuntos si están presentes
	if len(message.Attachments) > 0 {
		attachments := make([]brevo.SendSmtpEmailAttachment, len(message.Attachments))
		for i, attachment := range message.Attachments {
			// Codificar contenido en base64 si no está ya codificado
			content := base64.StdEncoding.EncodeToString(attachment.Content)

			attachments[i] = brevo.SendSmtpEmailAttachment{
				Name:    attachment.Name,
				Content: content,
			}
		}
		sendSmtpEmail.Attachment = attachments
	}

	return sendSmtpEmail
}

// NewBrevoEmailServiceFromEnv crea un servicio de email usando variables de entorno
func NewBrevoEmailServiceFromEnv(log logger.Logger) (*BrevoEmailService, error) {
	apiKey := os.Getenv(dto.EnvBrevoAPIKey)
	if apiKey == "" {
		return nil, errors.New(dto.ErrBrevoEnvAPIKeyRequired)
	}

	fromEmail := os.Getenv(dto.EnvBrevoFromEmail)
	if fromEmail == "" {
		fromEmail = dto.BrevoDefaultFromEmail // Email por defecto
	}

	fromName := os.Getenv(dto.EnvBrevoFromName)
	if fromName == "" {
		fromName = dto.BrevoDefaultFromName // Nombre por defecto
	}

	config := &interfaces.EmailServiceConfig{
		APIKey:    apiKey,
		FromEmail: fromEmail,
		FromName:  fromName,
	}

	return NewBrevoEmailService(config, log)
}
