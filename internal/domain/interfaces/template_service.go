package interfaces

// TemplateService define la interfaz para el servicio de plantillas
type TemplateService interface {
	// RenderWelcomeEmail renderiza la plantilla de email de bienvenida
	RenderWelcomeEmail(userName, password string) (string, error)

	// RenderPasswordResetEmail renderiza la plantilla de email para reset de contraseña
	RenderPasswordResetEmail(userName, resetLink string) (string, error)

	// RenderEmailValidation renderiza la plantilla de email para validación
	RenderEmailValidation(userName, validationLink string) (string, error)
}
