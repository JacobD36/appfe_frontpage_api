# Template Service - Clean Architecture

## Descripción

El servicio de templates ha sido refactorizado siguiendo los principios de **Clean Architecture** y **SOLID** para separar las responsabilidades y hacer el código más mantenible.

## Arquitectura

### Antes (❌ Violaciones):
- Templates HTML mezclados en `user_service.go`
- Violación del **Single Responsibility Principle (SRP)**
- Difícil de mantener y testear
- Lógica de presentación acoplada con lógica de negocio

### Después (✅ Clean Architecture):
- **Domain Layer**: Interface `TemplateService` en `internal/domain/interfaces/`
- **Adapter Layer**: Implementación `htmlTemplateService` en `internal/adapter/template/`
- **Usecase Layer**: Uso de la interfaz en `user_service.go`

## Estructura

```
internal/
├── domain/interfaces/
│   └── template_service.go          # Interface del dominio
├── adapter/template/
│   ├── html_template_service.go     # Implementación HTML
│   ├── templates.go                 # Constantes de templates
│   └── html_template_service_test.go # Tests unitarios
└── usecase/
    ├── user_service.go              # Usa la interfaz TemplateService
    └── dto/
        └── validation_messages.go   # Constantes de mensajes
```

## Beneficios

### 1. **Single Responsibility Principle (SRP)**
- `UserService`: Solo maneja lógica de negocio de usuarios
- `TemplateService`: Solo maneja renderizado de plantillas

### 2. **Open/Closed Principle (OCP)**
- Fácil agregar nuevos tipos de templates sin modificar código existente
- Extensible para diferentes formatos (HTML, texto plano, PDF)

### 3. **Dependency Inversion Principle (DIP)**
- `UserService` depende de la abstracción `TemplateService`
- No depende de implementaciones concretas

### 4. **Separation of Concerns**
- Presentación separada de lógica de negocio
- Templates fáciles de modificar por diseñadores

## Uso

### Templates disponibles:

```go
// Email de bienvenida
welcomeHTML, err := templateService.RenderWelcomeEmail("Juan Pérez", "password123")

// Reset de contraseña
resetHTML, err := templateService.RenderPasswordResetEmail("Juan Pérez", "https://app.com/reset?token=abc")

// Validación de email
validationHTML, err := templateService.RenderEmailValidation("Juan Pérez", "https://app.com/validate?token=xyz")
```

### Ejemplo en UserService:

```go
func (s *userService) Create(ctx context.Context, u *domain.User) error {
    // ... lógica de creación ...
    
    if s.messagingService != nil && s.templateService != nil {
        welcomeContent, err := s.templateService.RenderWelcomeEmail(u.Name, originalPassword)
        if err != nil {
            return uow.Commit() // No fallar por error de template
        }

        go func() {
            s.messagingService.SendEmail(ctx, u.Email, dto.WelcomeEmailSubject, welcomeContent)
        }()
    }
    
    return uow.Commit()
}
```

## Testing

Los templates tienen pruebas unitarias completas:

```bash
go test ./internal/adapter/template -v
```

## Extensibilidad

Para agregar nuevos tipos de templates:

1. **Actualizar la interfaz** en `template_service.go`:
```go
type TemplateService interface {
    RenderWelcomeEmail(userName, password string) (string, error)
    RenderPasswordResetEmail(userName, resetLink string) (string, error)
    RenderEmailValidation(userName, validationLink string) (string, error)
    // ➕ Nuevo método
    RenderNotificationEmail(userName, message string) (string, error)
}
```

2. **Implementar el método** en `html_template_service.go`:
```go
func (t *htmlTemplateService) RenderNotificationEmail(userName, message string) (string, error) {
    // Implementación...
}
```

3. **Agregar constantes** en `templates.go` y `validation_messages.go`

## Implementaciones alternativas futuras

La interfaz permite múltiples implementaciones:

- `htmlTemplateService`: Templates HTML (actual)
- `textTemplateService`: Templates texto plano
- `pdfTemplateService`: Templates PDF
- `multiLanguageTemplateService`: Templates multiidioma

## Patrones aplicados

- ✅ **Strategy Pattern**: Para diferentes tipos de rendering
- ✅ **Dependency Injection**: Templates inyectados en UserService  
- ✅ **Single Responsibility**: Cada clase tiene una responsabilidad
- ✅ **Interface Segregation**: Interfaces específicas y pequeñas
- ✅ **Clean Architecture**: Separación clara de capas
