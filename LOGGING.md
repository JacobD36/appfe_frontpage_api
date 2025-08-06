# Sistema de Logging Mejorado - APPFE Lima API

## 🚀 Características del Nuevo Sistema de Logging

El sistema de logging ha sido completamente mejorado para proporcionar información más clara y útil durante el desarrollo y producción.

### ✨ Mejoras Implementadas

1. **Logging Estructurado**: Utiliza formato JSON con campos estructurados
2. **Niveles de Log**: DEBUG, INFO, WARN, ERROR, FATAL
3. **Contexto Enriquecido**: Request IDs, User IDs, timestamps precisos
4. **Logging de Performance**: Duración de requests y operaciones de base de datos
5. **Logging de Autenticación**: Seguimiento de intentos de login y autenticación
6. **Logging de Operaciones de Negocio**: Seguimiento de operaciones importantes

### 🔧 Configuración

En tu archivo `.env`, configura el nivel de logging:

```bash
# Niveles disponibles: DEBUG, INFO, WARN, ERROR
LOG_LEVEL=INFO
```

### 📊 Tipos de Logs

#### 1. **Logs de HTTP Requests**
```json
{
  "timestamp": "2024-08-05 14:30:25.123",
  "level": "INFO",
  "msg": "HTTP Request",
  "request_id": "req_abc123",
  "method": "POST",
  "path": "/api/v1/users",
  "remote_addr": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "status_code": 201,
  "duration": "45ms"
}
```

#### 2. **Logs de Base de Datos**
```json
{
  "timestamp": "2024-08-05 14:30:25.456",
  "level": "DEBUG",
  "msg": "Database Operation Success",
  "request_id": "req_abc123",
  "operation": "CREATE",
  "table": "users",
  "duration": "12ms"
}
```

#### 3. **Logs de Autenticación**
```json
{
  "timestamp": "2024-08-05 14:30:25.789",
  "level": "INFO",
  "msg": "Authentication Success",
  "request_id": "req_abc123",
  "email": "user@example.com",
  "success": true,
  "reason": "valid_credentials"
}
```

#### 4. **Logs de Operaciones de Negocio**
```json
{
  "timestamp": "2024-08-05 14:30:26.012",
  "level": "INFO",
  "msg": "Business Operation",
  "request_id": "req_abc123",
  "operation": "user_creation",
  "user_id": "user_123",
  "email": "user@example.com",
  "role": "USER_ROLE"
}
```

### 🛠️ Uso en el Código

#### Logging Básico
```go
import "github.com/JacobD36/appfe_frontpage_api/pkg/logger"

// Logging básico
logger.Info(ctx, "Operation completed successfully")
logger.Warn(ctx, "Something might be wrong")
logger.LogError(ctx, "An error occurred")

// Logging con campos adicionales
logger.Info(ctx, "User created", 
    logger.String("user_id", userID),
    logger.String("email", email),
    logger.String("role", role),
)
```

#### Logging en Handlers
```go
func (h *UserHandler) Create(c echo.Context) error {
    ctx := c.Request().Context()
    
    logger.Info(ctx, "Creating new user", 
        logger.String("operation", "user_create"))
    
    // ... lógica del handler
    
    if err != nil {
        logger.LogError(ctx, "Failed to create user", 
            logger.String("email", input.Email),
            logger.Error("error", err))
        return Error(c, http.StatusInternalServerError, err.Error())
    }
    
    logger.Info(ctx, "User created successfully", 
        logger.String("user_id", user.ID))
    
    return Success(c, http.StatusCreated, "User created", user)
}
```

#### Logging de Operaciones Específicas
```go
// Logging de autenticación
logger.GetLogger().LogAuthentication(ctx, email, true, "valid_credentials")

// Logging de operaciones de base de datos
logger.GetLogger().LogDatabaseOperation(ctx, "SELECT", "users", duration, err)

// Logging de operaciones de negocio
details := map[string]interface{}{
    "email": user.Email,
    "role": user.Role,
}
logger.GetLogger().LogBusinessOperation(ctx, "user_creation", userID, details)
```

### 🔍 Monitoring y Debugging

#### Filtrar Logs por Nivel
```bash
# Solo errores
grep '"level":"ERROR"' app.log

# Solo warnings y errores
grep -E '"level":"(WARN|ERROR)"' app.log
```

#### Seguir una Request Específica
```bash
# Usar el request_id para seguir una request completa
grep 'req_abc123' app.log
```

#### Análisis de Performance
```bash
# Requests más lentas
grep '"duration"' app.log | sort -k duration -nr | head -10

# Operaciones de base de datos más lentas
grep '"msg":"Database Operation"' app.log | grep '"duration"' | sort -nr
```

### 📈 Beneficios del Nuevo Sistema

1. **Debugging Más Fácil**: Request IDs permiten seguir requests completas
2. **Mejor Monitoreo**: Métricas de performance y errores
3. **Análisis de Comportamiento**: Logs de autenticación y operaciones de negocio
4. **Formato Consistente**: JSON estructurado facilita el parsing automático
5. **Configuración Flexible**: Diferentes niveles según el entorno

### 🚦 Niveles de Log Recomendados

- **Desarrollo**: `DEBUG` - Máxima información
- **Testing**: `INFO` - Información general sin ruido excesivo
- **Staging**: `WARN` - Solo advertencias y errores
- **Producción**: `ERROR` - Solo errores críticos

### 📝 Ejemplos de Output

**Inicio de la aplicación:**
```json
{"timestamp":"2024-08-05 14:30:20.000","level":"INFO","msg":"Starting APPFE Lima API Server","port":":8080"}
{"timestamp":"2024-08-05 14:30:20.100","level":"INFO","msg":"Initializing database connection","driver":"POSTGRES"}
{"timestamp":"2024-08-05 14:30:20.200","level":"INFO","msg":"Connected to PostgreSQL database successfully","host":"localhost","port":5432,"database":"appfe_db"}
{"timestamp":"2024-08-05 14:30:20.300","level":"INFO","msg":"RSA keys loaded successfully"}
{"timestamp":"2024-08-05 14:30:20.400","level":"INFO","msg":"Services initialized successfully"}
{"timestamp":"2024-08-05 14:30:20.500","level":"INFO","msg":"Starting HTTP server","address":":8080"}
```

**Request completa de creación de usuario:**
```json
{"timestamp":"2024-08-05 14:30:25.000","level":"INFO","msg":"Request Started","request_id":"req_abc123","method":"POST","path":"/api/v1/users","remote_addr":"192.168.1.100"}
{"timestamp":"2024-08-05 14:30:25.100","level":"INFO","msg":"Creating new user","request_id":"req_abc123","operation":"user_create"}
{"timestamp":"2024-08-05 14:30:25.200","level":"DEBUG","msg":"Database Operation Success","request_id":"req_abc123","operation":"SELECT","table":"users","duration":"5ms"}
{"timestamp":"2024-08-05 14:30:25.300","level":"DEBUG","msg":"Database Operation Success","request_id":"req_abc123","operation":"INSERT","table":"users","duration":"8ms"}
{"timestamp":"2024-08-05 14:30:25.400","level":"INFO","msg":"User created successfully","request_id":"req_abc123","user_id":"user_456","email":"test@example.com","role":"USER_ROLE"}
{"timestamp":"2024-08-05 14:30:25.500","level":"INFO","msg":"Successful response","request_id":"req_abc123","status_code":201,"message":"Usuario creado exitosamente"}
{"timestamp":"2024-08-05 14:30:25.600","level":"INFO","msg":"HTTP Request","request_id":"req_abc123","method":"POST","path":"/api/v1/users","status_code":201,"duration":"600ms"}
```

Este sistema de logging mejorado te proporcionará una visibilidad completa del comportamiento de tu aplicación, facilitando tanto el desarrollo como el mantenimiento en producción.
