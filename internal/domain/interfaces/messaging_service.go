package interfaces

import (
	"context"
)

// MessagingService define la interfaz para el servicio de mensajería interno
// Este servicio será usado por otros servicios como inyección de dependencias
type MessagingService interface {
	// SendEmail envía un correo electrónico simple
	SendEmail(ctx context.Context, to, subject, htmlContent string) error
}
