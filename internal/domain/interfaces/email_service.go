package interfaces

import "context"

// EmailMessage representa un mensaje de correo electrónico
type EmailMessage struct {
	To          []string
	Subject     string
	HtmlContent string
	TextContent string
	From        EmailAddress
	ReplyTo     *EmailAddress
	Attachments []EmailAttachment
}

// EmailAddress representa una dirección de correo electrónico
type EmailAddress struct {
	Email string
	Name  string
}

// EmailAttachment representa un archivo adjunto
type EmailAttachment struct {
	Name    string
	Content []byte
	Type    string
}

// EmailServiceConfig contiene la configuración para el servicio de email
type EmailServiceConfig struct {
	APIKey    string
	FromEmail string
	FromName  string
}

// EmailService define la interfaz para el servicio de correo electrónico
type EmailService interface {
	// SendEmail envía un correo electrónico individual
	SendEmail(ctx context.Context, message *EmailMessage) error

	// SendBulkEmail envía correos electrónicos en lotes
	SendBulkEmail(ctx context.Context, messages []*EmailMessage) error

	// SendTemplate envía un correo usando una plantilla
	SendTemplate(ctx context.Context, templateID int64, to []string, templateData map[string]interface{}) error

	// ValidateEmail valida si una dirección de correo es válida
	ValidateEmail(email string) error

	// GetEmailStatus obtiene el estado de un correo enviado
	GetEmailStatus(ctx context.Context, messageID string) (*EmailStatus, error)
}

// EmailStatus representa el estado de un correo enviado
type EmailStatus struct {
	MessageID string
	Status    string
	SentAt    string
	Error     string
}
