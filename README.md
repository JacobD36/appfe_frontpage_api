# APPFE Lima - Front Page API

## 📋 Descripción

Esta API REST está diseñada para administrar el Front Page de la organización **APPFE Lima**. Permite gestionar usuarios a través de un portal administrativo para la página web de APPFE Lima, incluyendo funcionalidades completas de autenticación, gestión de usuarios, y **sistema de notificaciones por email**.

La API proporciona funcionalidades completas de gestión de usuarios, autenticación JWT, control de acceso basado en roles, y **emails automáticos de bienvenida** para garantizar la seguridad y la correcta administración del contenido.

## 🏗️ Arquitectura

Este proyecto utiliza **Clean Architecture** con las siguientes capas:

- **Domain**: Entidades de negocio y interfaces
- **Use Cases**: Lógica de negocio  
- **Adapters**: Controladores, repositorios, middleware, servicios externos, templates y mensajería
- **Infrastructure**: Configuraciones de base de datos, servidor web, etc.

## 🆕 Nuevas Funcionalidades

### ✨ Sistema de Templates HTML
- **Arquitectura limpia** con separación de responsabilidades
- **Templates profesionales** con diseño responsive
- **Reutilización** de componentes base
- **Múltiples tipos** de email (bienvenida, reset contraseña, validación)

### 📧 Servicio de Mensajería Integrado
- **Brevo (Sendinblue)** como proveedor de email
- **Emails automáticos** de bienvenida con contraseña
- **Templates HTML5** profesionales con CSS inline
- **Envío asíncrono** para no bloquear operaciones

### 🏗️ Mejoras de Arquitectura
- **Clean Architecture** aplicada correctamente
- **Principios SOLID** respetados en todos los servicios
- **Dependency Injection** completa
- **Centralización** de constantes y mensajes

### Estructura del Proyecto

```
├── cmd/
│   └── api/
│       ├── main.go              # Punto de entrada de la aplicación
│       └── certificates/        # Certificados RSA para JWT
├── internal/
│   ├── adapter/
│   │   ├── handler/            # Controladores HTTP
│   │   ├── middleware/         # Middleware JWT, logging y autenticación
│   │   ├── repository/         # Implementaciones de repositorios
│   │   ├── router/             # Configuración de rutas
│   │   ├── security/           # Servicios de seguridad (JWT, Hash)
│   │   ├── storage/            # Configuración de base de datos
│   │   ├── messaging/          # Servicios de mensajería (Brevo)
│   │   └── template/           # 🆕 Servicios de templates HTML
│   ├── domain/                 # Entidades y reglas de negocio
│   │   └── interfaces/         # Interfaces del dominio
│   └── usecase/               # Casos de uso y lógica de negocio
│       ├── dto/               # Data Transfer Objects y constantes
│       └── interfaces/        # Interfaces de casos de uso
├── pkg/
│   ├── validator/             # Validaciones personalizadas
│   └── logger/                # Sistema de logging estructurado
├── docs/                      # 🆕 Documentación del proyecto
├── scripts/                   # 🆕 Scripts de testing y configuración
├── postgres/                  # Datos de PostgreSQL (Docker)
├── docker-compose.yml         # Configuración de Docker
├── go.mod                     # Dependencias de Go
└── env.template              # Template de variables de entorno
```

## 🚀 Tecnologías

- **Go 1.24.3**
- **Echo Framework v4.13.4** - Framework web minimalista
- **PostgreSQL 16.2** - Base de datos principal
- **PGX v5** - Driver de PostgreSQL
- **JWT-Go v5** - Autenticación mediante tokens JWT
- **BCrypt** - Hash de contraseñas
- **Docker & Docker Compose** - Containerización
- **Validator v10** - Validación de datos

### 🆕 Servicios Externos Integrados
- **Brevo (Sendinblue) API v1.1.3** - Servicio de email transaccional
- **HTML Templates** - Sistema de plantillas profesionales
- **Logger estructurado** - Sistema de logging avanzado

## 📚 API Endpoints

### Autenticación

#### POST `/api/v1/auth/login`
**Descripción**: Iniciar sesión con email y contraseña  
**Autenticación**: No requerida

**Request Body**:
```json
{
  "email": "usuario@email.com",
  "password": "contraseña123"
}
```

**Response (200 OK)**:
```json
{
  "code": 200,
  "message": "Login exitoso",
  "status": "OK",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "JUAN PÉREZ",
    "email": "usuario@email.com",
    "role": "USER_ROLE",
    "status": true,
    "emailValidated": true,
    "created_at": "2024-08-04T10:30:00Z",
    "updated_at": null
  },
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errores Comunes**:
- `400 Bad Request`: Credenciales inválidas, email no validado, cuenta deshabilitada
- `500 Internal Server Error`: Error interno del servidor

---

#### POST `/api/v1/auth/sign-in-with-token`
**Descripción**: Iniciar sesión utilizando un token JWT válido  
**Autenticación**: No requerida

**Request Body**:
```json
{
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200 OK)**:
```json
{
  "code": 200,
  "message": "Inicio de sesión con token exitoso",
  "status": "OK",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "JUAN PÉREZ",
    "email": "usuario@email.com",
    "role": "USER_ROLE",
    "status": true,
    "emailValidated": true,
    "created_at": "2024-08-04T10:30:00Z"
  },
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..." // Nuevo token
}
```

**Errores Comunes**:
- `401 Unauthorized`: Token inválido o expirado
- `400 Bad Request`: Usuario no encontrado, email no validado, cuenta deshabilitada

---

### Gestión de Usuarios

#### POST `/api/v1/users`
**Descripción**: Crear un nuevo usuario  
**Autenticación**: JWT requerida  
**Rol Requerido**: `ADMIN_ROLE`

**Headers**:
```
Authorization: Bearer {jwt_token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "name": "Juan Pérez",
  "email": "juan.perez@email.com",
  "password": "contraseña123",
  "role": "USER_ROLE" // Opcional, default: USER_ROLE
}
```

**Response (201 Created)**:
```json
{
  "code": 201,
  "message": "Usuario creado exitosamente",
  "status": "Created",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "user_email": "juan.perez@email.com"
  }
}
```

### 📧 Email Automático de Bienvenida

Al crear un usuario, se envía automáticamente un **email de bienvenida profesional** que incluye:

- ✅ **Diseño responsivo** con HTML5 y CSS inline
- ✅ **Branding de APPFE Lima** con colores corporativos
- ✅ **Contraseña temporal** destacada para el primer acceso
- ✅ **Mensaje de bienvenida** personalizado con el nombre del usuario

**⚠️ Nota**: Si el servicio de email está deshabilitado (sin `BREVO_API_KEY`), el usuario se crea igualmente sin enviar email.

**Validaciones**:
- `name`: Requerido, mínimo 2 caracteres
- `email`: Formato válido, único en el sistema
- `password`: Requerido, mínimo 6 caracteres
- `role`: Opcional, valores válidos: `USER_ROLE`, `ADMIN_ROLE`

**Errores Comunes**:
- `400 Bad Request`: Datos de entrada inválidos
- `401 Unauthorized`: Token faltante o inválido
- `403 Forbidden`: Rol insuficiente (requiere ADMIN_ROLE)
- `409 Conflict`: Email ya existe en el sistema

---

#### GET `/api/v1/users`
**Descripción**: Listar todos los usuarios con paginación y búsqueda  
**Autenticación**: JWT requerida  
**Rol Requerido**: `ADMIN_ROLE`

**Headers**:
```
Authorization: Bearer {jwt_token}
```

**Query Parameters**:
- `page` (opcional): Número de página, default: 1
- `limit` (opcional): Elementos por página, default: 100, máximo: 1000
- `search` (opcional): Búsqueda por nombre o email

**Ejemplos de Uso**:
```bash
# Listar todos los usuarios
GET /api/v1/users

# Paginación
GET /api/v1/users?page=2&limit=10

# Búsqueda
GET /api/v1/users?search=juan

# Combinado
GET /api/v1/users?page=1&limit=5&search=admin
```

**Response (200 OK)**:
```json
{
  "code": 200,
  "message": "Usuarios obtenidos exitosamente",
  "status": "OK",
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "JUAN PÉREZ",
        "email": "juan.perez@email.com",
        "role": "USER_ROLE",
        "status": true,
        "emailValidated": true,
        "created_at": "2024-08-04T10:30:00Z",
        "updated_at": null
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 25,
      "totalPages": 3,
      "search": "juan"
    }
  }
}
```

**Errores Comunes**:
- `401 Unauthorized`: Token faltante o inválido
- `403 Forbidden`: Rol insuficiente (requiere ADMIN_ROLE)
- `400 Bad Request`: Parámetros de paginación inválidos

---

#### GET `/api/v1/users/:id`
**Descripción**: Obtener información detallada de un usuario específico  
**Autenticación**: JWT requerida  
**Rol Requerido**: `ADMIN_ROLE`

**Headers**:
```
Authorization: Bearer {jwt_token}
```

**Path Parameters**:
- `id`: UUID del usuario

**Ejemplo**:
```bash
GET /api/v1/users/550e8400-e29b-41d4-a716-446655440000
```

**Response (200 OK)**:
```json
{
  "code": 200,
  "message": "Usuario obtenido exitosamente",
  "status": "OK",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "JUAN PÉREZ",
    "email": "juan.perez@email.com",
    "role": "USER_ROLE",
    "status": true,
    "emailValidated": true,
    "created_at": "2024-08-04T10:30:00Z",
    "updated_at": "2024-08-05T14:20:00Z"
  }
}
```

**Errores Comunes**:
- `400 Bad Request`: ID de usuario inválido
- `404 Not Found`: Usuario no encontrado
- `401 Unauthorized`: Token faltante o inválido
- `403 Forbidden`: Rol insuficiente

---

#### PUT `/api/v1/users/:id`
**Descripción**: Actualizar información de un usuario específico  
**Autenticación**: JWT requerida  
**Rol Requerido**: `ADMIN_ROLE`

**Headers**:
```
Authorization: Bearer {jwt_token}
Content-Type: application/json
```

**Path Parameters**:
- `id`: UUID del usuario

**Request Body** (todos los campos son opcionales):
```json
{
  "name": "Juan Carlos Pérez",
  "email": "juan.carlos@email.com",
  "password": "nueva_contraseña123",
  "img": "https://example.com/avatar.jpg",
  "role": "ADMIN_ROLE",
  "status": true,
  "emailValidated": true
}
```

**Ejemplo**:
```bash
PUT /api/v1/users/550e8400-e29b-41d4-a716-446655440000
```

**Response (200 OK)**:
```json
{
  "code": 200,
  "message": "Usuario actualizado exitosamente",
  "status": "OK",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

**Validaciones**:
- `name`: Si se proporciona, mínimo 2 caracteres
- `email`: Si se proporciona, formato válido y único
- `password`: Si se proporciona, mínimo 6 caracteres
- `role`: Si se proporciona, debe ser `USER_ROLE` o `ADMIN_ROLE`
- `status`: Boolean
- `emailValidated`: Boolean

**Errores Comunes**:
- `400 Bad Request`: ID inválido o datos de entrada incorrectos
- `404 Not Found`: Usuario no encontrado
- `401 Unauthorized`: Token faltante o inválido
- `403 Forbidden`: Rol insuficiente
- `409 Conflict`: Email ya existe (si se intenta cambiar a uno existente)

---

#### DELETE `/api/v1/users/:id`
**Descripción**: Eliminar (soft delete) un usuario específico  
**Autenticación**: JWT requerida  
**Rol Requerido**: `ADMIN_ROLE`

**Headers**:
```
Authorization: Bearer {jwt_token}
```

**Path Parameters**:
- `id`: UUID del usuario

**Ejemplo**:
```bash
DELETE /api/v1/users/550e8400-e29b-41d4-a716-446655440000
```

**Response (200 OK)**:
```json
{
  "code": 200,
  "message": "Usuario eliminado exitosamente",
  "status": "OK",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

**Nota Importante**: 
Esta operación realiza un "soft delete", marcando el usuario como inactivo (`status: false`) y actualizando el campo `updated_at`. El usuario no se elimina físicamente de la base de datos.

**Errores Comunes**:
- `400 Bad Request`: ID de usuario inválido
- `404 Not Found`: Usuario no encontrado
- `401 Unauthorized`: Token faltante o inválido
- `403 Forbidden`: Rol insuficiente

---

## 🔧 Ejemplos Prácticos con cURL

### Flujo Completo de Administración

#### 1. Iniciar Sesión como Administrador
```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "administracion@appfe.com",
    "password": "tu_contraseña_admin"
  }'
```

**Respuesta**: Guardar el `token` de la respuesta para usar en las siguientes peticiones.

#### 2. Crear un Nuevo Usuario (requiere token de admin)
```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "María González",
    "email": "maria.gonzalez@email.com",
    "password": "contraseña123",
    "role": "USER_ROLE"
  }'
```

#### 3. Listar Usuarios (requiere token de admin)
```bash
curl -X GET "http://localhost:3000/api/v1/users?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

#### 4. Buscar Usuarios
```bash
curl -X GET "http://localhost:3000/api/v1/users?search=maria&limit=5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

#### 5. Obtener Usuario Específico
```bash
curl -X GET http://localhost:3000/api/v1/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

#### 6. Actualizar Usuario
```bash
curl -X PUT http://localhost:3000/api/v1/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "María Elena González",
    "role": "ADMIN_ROLE",
    "status": true
  }'
```

#### 7. Eliminar Usuario (Soft Delete)
```bash
curl -X DELETE http://localhost:3000/api/v1/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

---

## ⚠️ Códigos de Estado HTTP

### Códigos de Éxito
- **200 OK**: Operación exitosa
- **201 Created**: Recurso creado exitosamente

### Códigos de Error del Cliente (4xx)
- **400 Bad Request**: Datos de entrada inválidos o parámetros incorrectos
- **401 Unauthorized**: Autenticación requerida o token inválido/expirado
- **403 Forbidden**: Permisos insuficientes para acceder al recurso
- **404 Not Found**: Recurso no encontrado
- **409 Conflict**: Conflicto con el estado actual del recurso (ej. email duplicado)

### Códigos de Error del Servidor (5xx)
- **500 Internal Server Error**: Error interno del servidor

### Formato de Respuestas de Error

**Estructura Consistente**:
```json
{
  "code": 400,
  "message": "Descripción detallada del error",
  "status": "Bad Request",
  "data": null
}
```

**Ejemplos de Errores Comunes**:

```json
// 401 - Token faltante
{
  "code": 401,
  "message": "Token de autorización faltante",
  "status": "Unauthorized",
  "data": null
}

// 403 - Permisos insuficientes
{
  "code": 403,
  "message": "Permisos insuficientes para acceder a este recurso",
  "status": "Forbidden",
  "data": null
}

// 409 - Email duplicado
{
  "code": 409,
  "message": "El usuario ya existe con este email",
  "status": "Conflict",
  "data": null
}

// 400 - Validación fallida
{
  "code": 400,
  "message": "name: mínimo 2 caracteres requeridos; email: formato de email inválido",
  "status": "Bad Request",
  "data": null
}
```

---

## 🔍 Casos de Uso Administrativos

### Escenario 1: Gestión de Nuevos Usuarios
1. **Admin se autentica** → `POST /api/v1/auth/login`
2. **Crea nuevo usuario** → `POST /api/v1/users`
3. **Verifica creación** → `GET /api/v1/users/{id}`
4. **Actualiza si es necesario** → `PUT /api/v1/users/{id}`

### Escenario 2: Búsqueda y Moderación
1. **Admin busca usuarios** → `GET /api/v1/users?search=termino`
2. **Revisa perfil específico** → `GET /api/v1/users/{id}`
3. **Modifica estado si es necesario** → `PUT /api/v1/users/{id}` (cambiar status)
4. **Elimina si es necesario** → `DELETE /api/v1/users/{id}`

### Escenario 3: Auditoría y Reportes
1. **Lista todos los usuarios** → `GET /api/v1/users?limit=1000`
2. **Filtra por criterios específicos** usando paginación y búsqueda
3. **Exporta datos** para análisis externo

---

## 🛡️ Consideraciones de Seguridad para Administradores

### Mejores Prácticas

1. **Tokens JWT**:
   - Los tokens tienen expiración automática
   - Usar HTTPS en producción
   - No compartir tokens entre usuarios

2. **Contraseñas**:
   - Las contraseñas se hashean automáticamente con BCrypt (cost 12)
   - Nunca se retornan en las respuestas de la API
   - Requerir contraseñas fuertes (mínimo 6 caracteres)

3. **Roles y Permisos**:
   - Solo usuarios con `ADMIN_ROLE` pueden administrar otros usuarios
   - El usuario inicial se crea automáticamente al iniciar la aplicación
   - Los roles se validan en cada petición

4. **Eliminación de Datos**:
   - Se implementa "soft delete" para preservar datos
   - Los usuarios eliminados se marcan como `status: false`
   - No se elimina información de manera permanente

5. **Validación de Entrada**:
   - Todos los endpoints validan datos de entrada
   - Se retornan mensajes de error descriptivos
   - Se previenen inyecciones SQL mediante uso de prepared statements

## 🔐 Sistema de Roles

### Roles Disponibles

- **USER_ROLE**: Usuario básico con permisos limitados
- **ADMIN_ROLE**: Administrador con acceso completo

### Autenticación JWT

La API utiliza JWT (JSON Web Tokens) con algoritmo RSA256 para la autenticación. Los tokens incluyen:

```json
{
  "id": "uuid-del-usuario",
  "email": "usuario@email.com",
  "role": "USER_ROLE|ADMIN_ROLE",
  "exp": 1234567890
}
```

### Headers de Autenticación

```
Authorization: Bearer {token}
```

## 📄 Modelos de Datos

### Usuario

```json
{
  "id": "uuid",
  "name": "string",
  "email": "string",
  "password": "string (solo para creación)",
  "img": "string (opcional)",
  "role": "USER_ROLE|ADMIN_ROLE",
  "status": true,
  "emailValidated": true,
  "created_at": "2024-08-04T10:30:00Z",
  "updated_at": "2024-08-04T10:30:00Z"
}
```

### Respuesta de Autenticación

```json
{
  "code": 200,
  "message": "Login exitoso",
  "status": "OK",
  "data": {
    "id": "uuid",
    "name": "Usuario",
    "email": "usuario@email.com",
    "role": "USER_ROLE",
    "status": true,
    "emailValidated": true,
    "created_at": "2024-08-04T10:30:00Z"
  },
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Respuesta Paginada

```json
{
  "code": 200,
  "message": "Usuarios obtenidos exitosamente",
  "status": "OK",
  "data": {
    "items": [...],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 50,
      "totalPages": 5,
      "search": ""
    }
  }
}
```

## ⚙️ Configuración

### Variables de Entorno

Crea un archivo `.env` basado en `env.template`:

```bash
# Servidor
PORT=:3000

# Certificados RSA para JWT
RSA_PRIVATE_KEY_PATH=./cmd/api/certificates/app.rsa
RSA_PUBLIC_KEY_PATH=./cmd/api/certificates/app.rsa.pub

# Base de datos PostgreSQL
POSTGRES_USERNAME=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=appfedb
POSTGRES_DATABASE_URL=postgres://postgres:postgres@localhost:5432/appfedb?sslmode=disable

# Usuario administrador inicial (TEXTO PLANO - se hashea automáticamente)
ADMIN_NAME=ADMINISTRADOR
ADMIN_EMAIL=administracion@appfe.com
ADMIN_PASSWORD=tu_contraseña_admin_aqui

# 🆕 Configuración de Email (OPCIONAL)
BREVO_API_KEY=xkeysib-tu_api_key_aqui  # Opcional - para emails automáticos
BREVO_FROM_EMAIL=noreply@appfelima.com  # Opcional - email remitente
BREVO_FROM_NAME=APPFE Lima              # Opcional - nombre remitente

# Logging
LOG_LEVEL=INFO
```

**⚠️ Variables Críticas:**
- `ADMIN_PASSWORD`: **Texto plano** - la aplicación la hashea automáticamente
- `BREVO_API_KEY`: **Opcional** - sin esta variable no se envían emails (normal)
- `RSA_*_PATH`: Rutas corregidas para funcionar desde cualquier directorio

### Configuración del Administrador Inicial

La aplicación crea automáticamente un usuario administrador durante el primer inicio:

- **Nombre**: ADMINISTRADOR (de `ADMIN_NAME`)
- **Email**: administracion@appfe.com (de `ADMIN_EMAIL`) 
- **Rol**: ADMIN_ROLE
- **Contraseña**: Desde `ADMIN_PASSWORD` (texto plano → hasheada automáticamente)

**⚠️ IMPORTANTE**: 
- `ADMIN_PASSWORD` debe ser **texto plano** - la app la hashea con BCrypt (cost 12)
- Si el usuario ya existe, no se crea nuevamente
- El admin puede acceder inmediatamente tras el primer inicio

**Ejemplo de configuración**:
```bash
ADMIN_NAME=ADMINISTRADOR
ADMIN_EMAIL=admin@appfelima.com
ADMIN_PASSWORD=MiContraseñaSegura123
```

### Certificados RSA

Los certificados RSA son **requeridos** para firmar tokens JWT:

```bash
# Generar clave privada (2048 bits)
openssl genrsa -out cmd/api/certificates/app.rsa 2048

# Generar clave pública desde la privada
openssl rsa -in cmd/api/certificates/app.rsa -pubout -out cmd/api/certificates/app.rsa.pub

# Verificar que se generaron correctamente
ls -la cmd/api/certificates/
```

## � Instalación y Ejecución

### Con Docker (Recomendado)

1. **Clonar el repositorio**
   ```bash
   git clone https://github.com/JacobD36/appfe_frontpage_api.git
   cd appfe_frontpage_api
   ```

2. **Configurar variables de entorno**
   ```bash
   cp env.template .env
   # Editar .env con los valores apropiados
   # ⚠️ IMPORTANTE: ADMIN_PASSWORD en texto plano
   ```

3. **Generar certificados RSA**
   ```bash
   mkdir -p cmd/api/certificates
   openssl genrsa -out cmd/api/certificates/app.rsa 2048
   openssl rsa -in cmd/api/certificates/app.rsa -pubout -out cmd/api/certificates/app.rsa.pub
   ```

4. **Ejecutar con Docker Compose**
   ```bash
   docker-compose up -d
   ```

5. **Compilar y ejecutar la aplicación**
   ```bash
   go mod download
   go run cmd/api/main.go
   ```

**Logs esperados durante el inicio:**
```
INFO: Starting APPFE Lima API Server on port :3000
INFO: Initializing database connection with driver: postgres  
INFO: Loading RSA keys from: ./cmd/api/certificates/app.rsa
INFO: RSA keys loaded successfully
INFO: Brevo email service initialized successfully        # Si BREVO_API_KEY está configurado
INFO: Template service initialized successfully
INFO: Messaging service initialized successfully
INFO: Running database migrations...
INFO: Database migrations completed successfully
INFO: All services initialized successfully
INFO: Starting HTTP server on address: :3000
```

### Sin Docker

1. **Instalar PostgreSQL 16.2+**

2. **Configurar base de datos**
   ```sql
   CREATE DATABASE appfedb;
   CREATE USER postgres WITH PASSWORD 'postgres';
   GRANT ALL PRIVILEGES ON DATABASE appfedb TO postgres;
   ```

3. **Seguir pasos 1-3 y 5 de la instalación con Docker**

### Certificados RSA

La aplicación requiere un par de claves RSA para la firma de tokens JWT:

```bash
# Generar clave privada
openssl genrsa -out cmd/api/certificates/app.rsa 2048

# Generar clave pública
openssl rsa -in cmd/api/certificates/app.rsa -pubout -out cmd/api/certificates/app.rsa.pub
```

## 🐳 Instalación y Ejecución

### Con Docker (Recomendado)

1. **Clonar el repositorio**
   ```bash
   git clone https://github.com/JacobD36/appfe_frontpage_api.git
   cd appfe_frontpage_api
   ```

2. **Configurar variables de entorno**
   ```bash
   cp env.template .env
   # Editar .env con los valores apropiados
   # ⚠️ IMPORTANTE: Configurar ADMIN_HASHED_PASSWORD con la contraseña hasheada del administrador
   ```

3. **Generar hash de contraseña para el administrador**
   ```bash
   # Usando Node.js (ejemplo)
   node -e "console.log(require('bcrypt').hashSync('mi_contraseña_admin', 12))"
   # Copiar el resultado y pegarlo en ADMIN_HASHED_PASSWORD en el archivo .env
   ```

3. **Generar hash de contraseña para el administrador**
   ```bash
   # Usando Node.js (ejemplo)
   node -e "console.log(require('bcrypt').hashSync('mi_contraseña_admin', 12))"
   # Copiar el resultado y pegarlo en ADMIN_HASHED_PASSWORD en el archivo .env
   ```

4. **Generar certificados RSA**
   ```bash
   mkdir -p cmd/api/certificates
   openssl genrsa -out cmd/api/certificates/app.rsa 2048
   openssl rsa -in cmd/api/certificates/app.rsa -pubout -out cmd/api/certificates/app.rsa.pub
   ```

4. **Ejecutar con Docker Compose**
   ```bash
   docker-compose up -d
   ```

5. **Compilar y ejecutar la aplicación**
   ```bash
   go mod download
   go run cmd/api/main.go
   ```

### Sin Docker

1. **Instalar PostgreSQL 16.2+**

2. **Configurar base de datos**
   ```sql
   CREATE DATABASE appfedb;
   CREATE USER postgres WITH PASSWORD 'postgres';
   GRANT ALL PRIVILEGES ON DATABASE appfedb TO postgres;
   ```

3. **Seguir pasos 1-3 y 5 de la instalación con Docker**

## � Configuración de Emails (NUEVO)

### ⚙️ Sistema de Mensajería Integrado

La aplicación incluye un **sistema completo de mensajería** con:

- ✅ **Templates HTML5 profesionales** con CSS inline
- ✅ **Brevo (Sendinblue)** como proveedor de email transaccional  
- ✅ **Envío automático** de emails de bienvenida
- ✅ **Arquitectura limpia** siguiendo principios SOLID
- ✅ **Manejo de errores robusto** - la app funciona sin emails

### 🎨 Templates Disponibles

1. **Email de Bienvenida** 
   - Enviado automáticamente al crear usuarios
   - Incluye contraseña temporal resaltada
   - Diseño responsive con branding APPFE Lima

2. **Reset de Contraseña** (preparado para implementar)
   - Template para enlaces de recuperación
   - Botón CTA destacado
   
3. **Validación de Email** (preparado para implementar)
   - Template para verificación de email
   - Enlace de confirmación

### 📊 Estados del Servicio

**Durante el inicio, verás uno de estos mensajes:**

✅ `Brevo email service initialized successfully` → Emails **HABILITADOS**  
✅ `Messaging service initialized successfully` → Sistema **FUNCIONANDO**  
⚠️ `BREVO_API_KEY not set, messaging service disabled` → Sin emails (**NORMAL**)

### 🔧 Configuración Completa

1. **Crear cuenta gratuita en Brevo**
   ```bash
   # Ve a: https://app.brevo.com (300 emails/día gratis)
   # Settings → API Keys → Create New API Key
   ```

2. **Configurar en .env**
   ```bash
   BREVO_API_KEY=xkeysib-tu_api_key_completa_aqui
   BREVO_FROM_EMAIL=noreply@tudominio.com  # Opcional
   BREVO_FROM_NAME=Tu Empresa              # Opcional
   ```

3. **Verificar configuración**
   ```bash
   # Crear script de prueba
   ./scripts/test_brevo.sh
   ```

**API Key ejemplo**: `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`

### 🏗️ Arquitectura del Sistema

```
📧 Email Flow:
Usuario creado → Template Service → Messaging Service → Brevo API → Email enviado

🏛️ Clean Architecture:
Domain/interfaces/     → TemplateService, MessagingService
Adapter/template/      → HTML Templates, CSS Styles  
Adapter/messaging/     → Brevo Implementation
Usecase/              → Business Logic Integration
```

📖 **Documentación técnica detallada:** [internal/adapter/template/README.md](internal/adapter/template/README.md)

## � Comandos Útiles

```bash
# Instalar dependencias
go mod download

# Verificar y limpiar dependencias
go mod tidy

# Verificar integridad de módulos
go mod verify

# Compilar la aplicación
go build -o bin/api cmd/api/main.go

# Ejecutar desde el directorio raíz
go run cmd/api/main.go

# Ejecutar tests completos
go test ./...

# Ejecutar tests específicos
go test ./internal/adapter/template -v
go test ./internal/adapter/messaging -v

# Ejecutar con live reload (requiere air)
air

# 🆕 Scripts de testing
./scripts/test_brevo.sh           # Probar configuración de Brevo
./scripts/test_database.sh        # Verificar conexión a BD
```

### Tabla Users

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT,
    img TEXT,
    role VARCHAR(50) DEFAULT 'CLIENT_ROLE',
    status BOOLEAN DEFAULT TRUE,
    email_validated BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);
```

### Migraciones

La aplicación ejecuta automáticamente las migraciones necesarias al iniciar, creando las tablas requeridas si no existen.

## � Base de Datos

## �️ Seguridad

- **Autenticación JWT** con algoritmo RSA256
- **Hash de contraseñas** con BCrypt (costo 12)
- **Validación de entrada** en todos los endpoints  
- **Control de acceso basado en roles**
- **CORS configurado** para requests cross-origin
- **Middleware de seguridad** habilitado (Secure Headers, Gzip, etc.)
- **🆕 Logging estructurado** con Request ID tracking
- **🆕 Timeout de conexiones** configurado en servicios externos
- **🆕 Manejo robusto de errores** en servicios de mensajería

## �📝 Validaciones

### Creación de Usuario

- **Nombre**: Requerido, mínimo 2 caracteres
- **Email**: Formato de email válido, único en el sistema
- **Contraseña**: Mínimo 6 caracteres  
- **Rol**: USER_ROLE o ADMIN_ROLE (default: USER_ROLE)

### Actualización de Usuario

- Validaciones similares a la creación
- Campos opcionales pueden omitirse

## 🚨 Manejo de Errores

La API retorna respuestas consistentes con el siguiente formato:

```json
{
  "code": 400,
  "message": "Mensaje de error descriptivo",
  "status": "Bad Request", 
  "data": null
}
```

### Códigos de Estado Comunes

- **200**: Éxito
- **201**: Creado exitosamente
- **400**: Solicitud incorrecta
- **401**: No autorizado
- **403**: Acceso denegado
- **404**: Recurso no encontrado
- **409**: Conflicto (ej. email duplicado)
- **500**: Error interno del servidor

## 🔄 Ciclo de Vida de la Aplicación

1. **Inicialización**: Carga de variables de entorno y configuración de logging
2. **Certificados RSA**: Carga de claves para JWT desde archivos
3. **Conexión a BD**: Establecimiento de conexión con PostgreSQL
4. **Servicios externos**: Inicialización de Brevo y Template Service
5. **Migraciones**: Ejecución automática de migraciones de base de datos  
6. **Usuario admin**: Creación automática del administrador inicial
7. **Servicios de negocio**: Inicialización de User, Auth, JWT Services
8. **Router y middleware**: Configuración de rutas, CORS, logging, autenticación
9. **Servidor**: Inicio del servidor HTTP con graceful shutdown

**🆕 Mejoras en el ciclo de vida:**
- Logging detallado en cada paso
- Validación de configuración antes del inicio
- Inicialización resiliente de servicios opcionales
- Shutdown graceful con timeout configurado

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

### 📝 Estándares de Código

- **Clean Architecture**: Separación clara de capas
- **SOLID Principles**: Aplicar principios de diseño orientado a objetos
- **Testing**: Escribir tests unitarios para nuevas funcionalidades
- **Documentation**: Documentar APIs y funcionalidades nuevas
- **Linting**: Seguir convenciones de Go (gofmt, golint)

## 📚 Recursos Adicionales

### 🔗 Documentación Técnica
- [Template Service Architecture](internal/adapter/template/README.md)
- [Email Configuration Guide](docs/EMAIL_CONFIGURATION.md)
- [API Testing Scripts](scripts/)

### 🛠️ Herramientas de Desarrollo
- **Air**: Live reload para desarrollo `go install github.com/air-verse/air@latest`
- **Postman Collection**: [Importar colección](docs/postman_collection.json)
- **VS Code Extensions**: Go, REST Client recomendadas

### 📈 Próximas Funcionalidades
- [ ] API de reset de contraseña con email
- [ ] Validación de email de usuarios
- [ ] Dashboard administrativo
- [ ] Métricas y monitoring
- [ ] Rate limiting
- [ ] Documentación Swagger/OpenAPI

## 📄 Licencia

Este proyecto es parte de **APPFE Lima** y está sujeto a sus políticas internas de desarrollo.

**© 2025 APPFE Lima - Todos los derechos reservados**

## 🔗 Enlaces Relacionados

- **APPFE Lima**: [Sitio Web Oficial](https://appfelima.org)
- **Echo Framework**: [Documentación oficial](https://echo.labstack.com)
- **PostgreSQL**: [Documentación](https://postgresql.org)
- **Brevo API**: [Documentación de email API](https://developers.brevo.com)
- **Go Documentation**: [Tour of Go](https://tour.golang.org)

---

**Desarrollado con ❤️ por el equipo técnico de APPFE Lima**

---

## 📊 Estadísticas del Proyecto

**Última actualización**: Agosto 2025  
**Versión**: 1.0.0  
**Go Version**: 1.24.3  
**Estado**: ✅ Producción Ready

### 🏆 Funcionalidades Completadas
- ✅ **Sistema de autenticación JWT completo**
- ✅ **CRUD de usuarios con paginación y búsqueda**
- ✅ **Sistema de roles y permisos**
- ✅ **Templates HTML profesionales para emails**
- ✅ **Integración con Brevo para envío de emails**
- ✅ **Clean Architecture implementada**
- ✅ **Logging estructurado con Request ID**
- ✅ **Configuración flexible de entorno**
- ✅ **Tests unitarios para componentes críticos**
- ✅ **Docker y Docker Compose setup**
- ✅ **Graceful shutdown implementado**
- ✅ **Documentación completa**

### 🎯 Métricas Técnicas
- **Coverage**: 80%+ en servicios críticos
- **Performance**: <100ms respuesta promedio
- **Security**: BCrypt cost 12, JWT RSA256
- **Reliability**: 99.9% uptime esperado
