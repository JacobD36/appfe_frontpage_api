package template

import (
	"fmt"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/dto"
)

type htmlTemplateService struct{}

func NewHTMLTemplateService() interfaces.TemplateService {
	return &htmlTemplateService{}
}

func (t *htmlTemplateService) RenderWelcomeEmail(userName, password string) (string, error) {
	if userName == "" || password == "" {
		return "", fmt.Errorf(dto.ErrTemplateRenderFailed, fmt.Errorf("userName and password are required"))
	}

	title := "Bienvenido a APPFE Lima"
	header := "Bienvenido(a) a APPFE Lima"
	content := fmt.Sprintf(welcomeContentTemplate, userName, password)

	return fmt.Sprintf(baseHTMLTemplate, title, header, content), nil
}

func (t *htmlTemplateService) RenderPasswordResetEmail(userName, resetLink string) (string, error) {
	if userName == "" || resetLink == "" {
		return "", fmt.Errorf(dto.ErrTemplateRenderFailed, fmt.Errorf("userName and resetLink are required"))
	}

	title := "Restablecer Contraseña - APPFE Lima"
	header := "Restablecer Contraseña"
	content := fmt.Sprintf(passwordResetContentTemplate, userName, resetLink)

	return fmt.Sprintf(baseHTMLTemplate, title, header, content), nil
}

func (t *htmlTemplateService) RenderEmailValidation(userName, validationLink string) (string, error) {
	if userName == "" || validationLink == "" {
		return "", fmt.Errorf(dto.ErrTemplateRenderFailed, fmt.Errorf("userName and validationLink are required"))
	}

	title := "Verificar Email - APPFE Lima"
	header := "Verificación de Email"
	content := fmt.Sprintf(emailValidationContentTemplate, userName, validationLink)

	return fmt.Sprintf(baseHTMLTemplate, title, header, content), nil
}
