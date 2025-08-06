package template

import (
	"strings"
	"testing"
)

func TestHTMLTemplateService_RenderWelcomeEmail(t *testing.T) {
	service := NewHTMLTemplateService()

	tests := []struct {
		name     string
		userName string
		password string
		wantErr  bool
	}{
		{
			name:     "valid parameters",
			userName: "John Doe",
			password: "tempPassword123",
			wantErr:  false,
		},
		{
			name:     "empty userName",
			userName: "",
			password: "tempPassword123",
			wantErr:  true,
		},
		{
			name:     "empty password",
			userName: "John Doe",
			password: "",
			wantErr:  true,
		},
		{
			name:     "both empty",
			userName: "",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.RenderWelcomeEmail(tt.userName, tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("RenderWelcomeEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verificar que el resultado contiene elementos esperados
				if !strings.Contains(result, "DOCTYPE html") {
					t.Error("Expected HTML5 doctype")
				}
				if !strings.Contains(result, tt.userName) {
					t.Error("Expected userName to be in template")
				}
				if !strings.Contains(result, tt.password) {
					t.Error("Expected password to be in template")
				}
				if !strings.Contains(result, "APPFE Lima") {
					t.Error("Expected company name in template")
				}
			}
		})
	}
}

func TestHTMLTemplateService_RenderPasswordResetEmail(t *testing.T) {
	service := NewHTMLTemplateService()

	tests := []struct {
		name      string
		userName  string
		resetLink string
		wantErr   bool
	}{
		{
			name:      "valid parameters",
			userName:  "John Doe",
			resetLink: "https://appfe.com/reset?token=abc123",
			wantErr:   false,
		},
		{
			name:      "empty userName",
			userName:  "",
			resetLink: "https://appfe.com/reset?token=abc123",
			wantErr:   true,
		},
		{
			name:      "empty resetLink",
			userName:  "John Doe",
			resetLink: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.RenderPasswordResetEmail(tt.userName, tt.resetLink)

			if (err != nil) != tt.wantErr {
				t.Errorf("RenderPasswordResetEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !strings.Contains(result, "DOCTYPE html") {
					t.Error("Expected HTML5 doctype")
				}
				if !strings.Contains(result, tt.userName) {
					t.Error("Expected userName to be in template")
				}
				if !strings.Contains(result, tt.resetLink) {
					t.Error("Expected resetLink to be in template")
				}
				if !strings.Contains(result, "Restablecer Contraseña") {
					t.Error("Expected password reset text in template")
				}
			}
		})
	}
}

func TestHTMLTemplateService_RenderEmailValidation(t *testing.T) {
	service := NewHTMLTemplateService()

	tests := []struct {
		name           string
		userName       string
		validationLink string
		wantErr        bool
	}{
		{
			name:           "valid parameters",
			userName:       "Jane Doe",
			validationLink: "https://appfe.com/validate?token=xyz789",
			wantErr:        false,
		},
		{
			name:           "empty userName",
			userName:       "",
			validationLink: "https://appfe.com/validate?token=xyz789",
			wantErr:        true,
		},
		{
			name:           "empty validationLink",
			userName:       "Jane Doe",
			validationLink: "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.RenderEmailValidation(tt.userName, tt.validationLink)

			if (err != nil) != tt.wantErr {
				t.Errorf("RenderEmailValidation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !strings.Contains(result, "DOCTYPE html") {
					t.Error("Expected HTML5 doctype")
				}
				if !strings.Contains(result, tt.userName) {
					t.Error("Expected userName to be in template")
				}
				if !strings.Contains(result, tt.validationLink) {
					t.Error("Expected validationLink to be in template")
				}
				if !strings.Contains(result, "Verificar Email") {
					t.Error("Expected email validation text in template")
				}
			}
		})
	}
}
