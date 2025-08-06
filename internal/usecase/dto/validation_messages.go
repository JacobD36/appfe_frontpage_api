package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Constantes para mensajes de validación y errores del sistema
const (
	// Mensajes de validación de contraseña
	ErrPasswordMinLength   = "la contraseña debe tener al menos 8 caracteres"
	ErrPasswordComplexity  = "la contraseña debe contener al menos una mayúscula, una minúscula, un número y un carácter especial"
	ErrPasswordsDoNotMatch = "las contraseñas no coinciden"

	// Mensajes de validación de entrada
	ErrInvalidInput   = "entrada inválida"
	ErrFieldRequired  = "el campo %s es obligatorio"
	ErrFieldMinLength = "el campo %s debe tener al menos %s caracteres"
	ErrInvalidEmail   = "el correo electrónico no es válido"
	ErrFieldInvalid   = "el campo %s no es válido"
	ErrEmptyPassword  = "el password no puede estar vacío"

	// Mensajes de usuario
	ErrUserAlreadyExists     = "El correo ya está registrado"
	ErrUserNotFound          = "Usuario no encontrado"
	ErrUserCreatedSuccess    = "Usuario creado exitosamente"
	ErrUserUpdatedSuccess    = "Usuario actualizado exitosamente"
	ErrUserDeletedSuccess    = "Usuario eliminado exitosamente"
	ErrUsersRetrievedSuccess = "Usuarios obtenidos exitosamente"
	ErrUserRetrievedSuccess  = "Usuario obtenido exitosamente"
	ErrInvalidUserID         = "ID de usuario inválido"

	// Mensajes de respuesta de datos de usuario
	UserDataLabel  = "user"
	UserIdLabel    = "id"
	UserEmailLabel = "email"

	// Mensajes de autenticación
	ErrInvalidCredentials    = "credenciales inválidas"
	ErrEmailNotValidated     = "el correo electrónico no ha sido validado"
	ErrAccountDisabled       = "la cuenta está deshabilitada"
	ErrLoginSuccess          = "inicio de sesión exitoso"
	ErrTokenGenerationFailed = "error al generar el token de acceso"

	// Mensajes de BcryptHasher
	ErrPasswordTooLong      = "la contraseña no puede exceder 72 caracteres"
	ErrHashEmpty            = "el hash de la contraseña no puede estar vacío"
	ErrHashGeneration       = "error al generar hash: %w"
	ErrPasswordIncorrect    = "contraseña incorrecta"
	ErrPasswordVerification = "error al verificar contraseña: %w"

	// Mensajes de servidor
	ErrInternalServer = "error interno del servidor"
	ErrNoRowsFound    = "no rows in result set"

	// Mensajes de sistema/main
	ErrLoadingEnvFile          = "Error loading .env file"
	ErrRSAKeysNotSet           = "RSA key paths not set in environment variables"
	ErrFailedLoadRSAKeys       = "Failed to load RSA keys: %v"
	ErrMigrationFailed         = "Migration failed: %v"
	ErrServerError             = "Server error: %v"
	ErrForcedShutdown          = "Forced shutdown: %v"
	ErrUnitOfWorkFactory       = "Failed to create UnitOfWork factory"
	MsgShuttingDownServer      = "Shutting down server..."
	MsgServerStoppedGracefully = "Server stopped gracefully"

	// Mensajes de base de datos/storage
	ErrUnsupportedDriver    = "Unsupported database driver: %s"
	ErrDatabaseURLNotSet    = "POSTGRES_DATABASE_URL environment variable is not set"
	ErrFailedParseConfig    = "Failed to parse config: %v"
	ErrUnableCreatePool     = "Unable to create connection pool: %v"
	ErrDriverNotImplemented = "Driver not implemented: %s"
	MsgConnectedToDatabase  = "Connected to PostgreSQL database successfully"

	// Mensajes de transacciones/Unit of Work
	ErrTransactionAlreadyCommitted  = "la transacción ya ha sido confirmada"
	ErrTransactionAlreadyRolledBack = "la transacción ya ha sido revertida"

	// Mensajes de paginación
	ErrInvalidPaginationPage   = "el número de página debe ser mayor a 0"
	ErrInvalidPaginationLimit  = "el límite debe estar entre 1 y 100"
	ErrUsersRetrievedPaginated = "Usuarios obtenidos exitosamente con paginación"
	ErrUsersSearchSuccess      = "Búsqueda de usuarios realizada exitosamente"

	// Mensajes de JWT y autorización
	ErrTokenMissing            = "token de autorización requerido"
	ErrInvalidTokenFormat      = "formato de token inválido. Use: Bearer <token>"
	ErrInvalidToken            = "token inválido o expirado"
	ErrInsufficientPermissions = "permisos insuficientes para acceder a este recurso"
	ErrTokenExpired            = "el token ha expirado"
	ErrUnauthorizedAccess      = "acceso no autorizado"
	ErrTokenSignInSuccess      = "autenticación con token exitosa"
	ErrUserNotFoundForToken    = "usuario no encontrado para el token proporcionado"

	// Mensajes de logging para el sistema
	MsgStartingAPIServer        = "Starting APPFE Lima API Server"
	MsgInitializingDBConnection = "Initializing database connection"
	MsgLoadingRSAKeys           = "Loading RSA keys"
	MsgRSAKeysLoadedSuccess     = "RSA keys loaded successfully"
	MsgRunningDBMigrations      = "Running database migrations"
	MsgDBMigrationsCompleted    = "Database migrations completed successfully"
	MsgServicesInitialized      = "Services initialized successfully"
	MsgStartingHTTPServer       = "Starting HTTP server"
	MsgShutdownSignalReceived   = "Shutdown signal received, starting graceful shutdown"

	// Mensajes de logging para handlers
	MsgCreatingNewUser           = "Creating new user"
	MsgUserCreatedSuccessfully   = "User created successfully"
	MsgInvalidInputForUser       = "Invalid input for user creation"
	MsgUserValidationFailed      = "User creation validation failed"
	MsgErrorCheckingExistingUser = "Error checking existing user"
	MsgAttemptCreateExistingUser = "Attempt to create user with existing email"
	MsgFailedToCreateUser        = "Failed to create user"

	// Mensajes de logging para respuestas
	MsgSuccessfulResponse      = "Successful response"
	MsgSuccessfulLoginResponse = "Successful login response"
	MsgErrorResponse           = "Error response"

	// Mensajes de logging para base de datos
	MsgInitializingPostgresConn = "Initializing PostgreSQL connection"
	MsgParsingDBConfig          = "Parsing database configuration"
	MsgCreatingDBPool           = "Creating database connection pool"
	MsgConnectedToDBSuccess     = "Connected to PostgreSQL database successfully"

	// Mensajes de logging para middleware y requests
	MsgRequestStarted        = "Request Started"
	MsgRequestError          = "Request Error"
	MsgHTTPRequest           = "HTTP Request"
	MsgDatabaseOpSuccess     = "Database Operation Success"
	MsgDatabaseOpFailed      = "Database Operation Failed"
	MsgAuthenticationSuccess = "Authentication Success"
	MsgAuthenticationFailed  = "Authentication Failed"
	MsgBusinessOperation     = "Business Operation"

	// Mensajes de error para mensajería
	ErrMessagingRecipientRequired  = "el destinatario es obligatorio"
	ErrMessagingSubjectRequired    = "el asunto es obligatorio"
	ErrMessagingContentRequired    = "el contenido es obligatorio"
	ErrMessagingInvalidEmailFormat = "formato de correo electrónico inválido: %w"
	ErrMessagingFailedToSend       = "error al enviar correo: %w"

	// Mensajes de logging para mensajería
	MsgMessagingSendingEmail      = "Sending email"
	MsgMessagingEmailSentSuccess  = "Email sent successfully"
	MsgMessagingFailedToSendEmail = "Failed to send email"

	// Mensajes para inicialización de servicios en main.go
	ErrMessagingServiceInitFailed  = "Failed to initialize email service"
	MsgMessagingServiceInitSuccess = "Messaging service initialized successfully"
	MsgMessagingServiceDisabled    = "BREVO_API_KEY not set, messaging service disabled"

	// Constantes para BrevoEmailService
	ErrBrevoConfigRequired            = "config is required"
	ErrBrevoAPIKeyRequired            = "API key is required"
	ErrBrevoFromEmailRequired         = "from email is required"
	ErrBrevoInvalidFromEmailFormat    = "invalid from email format: %w"
	ErrBrevoFailedValidateConfig      = "failed to validate Brevo configuration: %w"
	ErrBrevoInvalidAPIKeyOrConnection = "invalid API key or connection error: %w"
	ErrBrevoMessageRequired           = "message is required"
	ErrBrevoInvalidEmailMessage       = "invalid email message: %w"
	ErrBrevoFailedSendEmail           = "failed to send email: %w"
	ErrBrevoAtLeastOneMessageRequired = "at least one message is required"
	ErrBrevoTooManyMessages           = "too many messages in bulk request (max 50)"
	ErrBrevoFailedSendMultipleMsg     = "failed to send %d messages: %s"
	ErrBrevoValidTemplateIDRequired   = "valid template ID is required"
	ErrBrevoAtLeastOneRecipient       = "at least one recipient is required"
	ErrBrevoInvalidEmailAddress       = "invalid email address %s: %w"
	ErrBrevoFailedSendTemplate        = "failed to send template email: %w"
	ErrBrevoEmailRequired             = "email is required"
	ErrBrevoInvalidEmailFormat        = "invalid email format: %w"
	ErrBrevoMessageIDRequired         = "message ID is required"
	ErrBrevoSubjectRequired           = "subject is required"
	ErrBrevoContentRequired           = "either HTML content or text content is required"
	ErrBrevoEnvAPIKeyRequired         = "BREVO_API_KEY environment variable is required"

	// Constantes para logging de BrevoEmailService
	MsgBrevoServiceInitialized      = "Brevo email service initialized successfully"
	MsgBrevoFailedValidateConfig    = "Failed to validate Brevo API configuration"
	MsgBrevoInvalidEmailMessage     = "Invalid email message"
	MsgBrevoFailedSendEmail         = "Failed to send email via Brevo"
	MsgBrevoEmailSentSuccess        = "Email sent successfully via Brevo"
	MsgBrevoBulkEmailCompleted      = "Bulk email sending completed"
	MsgBrevoFailedSendTemplate      = "Failed to send template email via Brevo"
	MsgBrevoTemplateSentSuccess     = "Template email sent successfully via Brevo"
	MsgBrevoGetStatusNotImplemented = "GetEmailStatus is not fully implemented for Brevo"

	// Constantes para valores por defecto y configuración
	BrevoAPIKeyHeaderName   = "api-key"
	BrevoDefaultFromEmail   = "noreply@appfe.com"
	BrevoDefaultFromName    = "APPFE"
	BrevoDefaultEmailStatus = "sent"
	BrevoWebhookNote        = "Brevo requires webhooks for real-time status tracking"
	BrevoBulkMessageFormat  = "message %d: %s"

	// Constantes para variables de entorno
	EnvBrevoAPIKey    = "BREVO_API_KEY"
	EnvBrevoFromEmail = "BREVO_FROM_EMAIL"
	EnvBrevoFromName  = "BREVO_FROM_NAME"

	// Constantes para campos de logging
	LogFieldFromEmail  = "from_email"
	LogFieldFromName   = "from_name"
	LogFieldError      = "error"
	LogFieldSubject    = "subject"
	LogFieldTo         = "to"
	LogFieldMessageID  = "message_id"
	LogFieldTemplateID = "template_id"
	LogFieldTotalMsgs  = "total_messages"
	LogFieldSuccessful = "successful"
	LogFieldFailed     = "failed"
	LogFieldNote       = "note"

	// Mensajes de Template Service
	MsgTemplateServiceInitialized = "template service initialized successfully"
	ErrTemplateRenderFailed       = "error rendering template: %w"
	MsgWelcomeEmailTemplateUsed   = "welcome email template rendered"
	WelcomeEmailSubject           = "¡Bienvenido a APPFE Lima!"
	PasswordResetEmailSubject     = "Restablecer Contraseña - APPFE Lima"
	EmailValidationSubject        = "Verificar Email - APPFE Lima"
	ErrTemplateNotFound           = "template not found: %s"
	ErrTemplateParamsRequired     = "template parameters are required"
)

func TranslateValidationErrors(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		msg := ""
		for i, fieldErr := range errs {
			if i > 0 {
				msg += ". "
			}
			switch fieldErr.Tag() {
			case "required":
				msg += fmt.Sprintf(ErrFieldRequired, fieldErr.Field())
			case "min":
				msg += fmt.Sprintf(ErrFieldMinLength, fieldErr.Field(), fieldErr.Param())
			case "email":
				msg += ErrInvalidEmail
			default:
				msg += fmt.Sprintf(ErrFieldInvalid, fieldErr.Field())
			}
		}
		return msg
	}

	return err.Error()
}
