# APPFE Lima - Front Page API

## 📋 Descripción

Esta API REST está diseñada para administrar el Front Page de la organización **APPFE Lima**. Permite gestionar usuarios a través de un portal administrativo para la página web de APPFE Lima, donde los usuarios pueden modificar el contenido del Front Page según su rol asignado.

La API proporciona funcionalidades completas de gestión de usuarios, autenticación JWT, y control de acceso basado en roles para garantizar la seguridad y la correcta administración del contenido.

## 🏗️ Arquitectura

Este proyecto utiliza **Clean Architecture** con las siguientes capas:

- **Domain**: Entidades de negocio y interfaces
- **Use Cases**: Lógica de negocio
- **Adapters**: Controladores, repositorios, middleware y servicios externos
- **Infrastructure**: Configuraciones de base de datos, servidor web, etc.

### Estructura del Proyecto

```
├── cmd/
│   └── api/
│       ├── main.go              # Punto de entrada de la aplicación
│       └── certificates/        # Certificados RSA para JWT
├── internal/
│   ├── adapter/
│   │   ├── handler/            # Controladores HTTP
│   │   ├── middleware/         # Middleware JWT y autenticación
│   │   ├── repository/         # Implementaciones de repositorios
│   │   ├── router/             # Configuración de rutas
│   │   ├── security/           # Servicios de seguridad (JWT, Hash)
│   │   └── storage/            # Configuración de base de datos
│   ├── domain/                 # Entidades y reglas de negocio
│   │   └── interfaces/         # Interfaces del dominio
│   └── usecase/               # Casos de uso y lógica de negocio
│       ├── dto/               # Data Transfer Objects
│       └── interfaces/        # Interfaces de casos de uso
├── pkg/
│   └── validator/             # Validaciones personalizadas
├── postgres/                  # Datos de PostgreSQL (Docker)
├── docker-compose.yml         # Configuración de Docker
├── go.mod                     # Dependencias de Go
└── env.template              # Template de variables de entorno
```

## 🚀 Tecnologías

- **Go 1.24.3**
- **Echo Framework** - Framework web minimalista
- **PostgreSQL 16.2** - Base de datos principal
- **PGX v5** - Driver de PostgreSQL
- **JWT-Go v5** - Autenticación mediante tokens JWT
- **BCrypt** - Hash de contraseñas
- **Docker & Docker Compose** - Containerización
- **Validator v10** - Validación de datos

## 📚 API Endpoints

### Autenticación

| Método | Endpoint | Descripción | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/auth/login` | Iniciar sesión | No |
| POST | `/api/v1/auth/sign-in-with-token` | Iniciar sesión con token | No |

### Gestión de Usuarios

| Método | Endpoint | Descripción | Auth | Rol Requerido |
|--------|----------|-------------|------|---------------|
| POST | `/api/v1/users` | Crear usuario | No | - |
| GET | `/api/v1/users` | Listar usuarios | JWT | ADMIN_ROLE |
| GET | `/api/v1/users/:id` | Obtener usuario por ID | JWT | ADMIN_ROLE |
| PUT | `/api/v1/users/:id` | Actualizar usuario | JWT | ADMIN_ROLE |
| DELETE | `/api/v1/users/:id` | Eliminar usuario | JWT | ADMIN_ROLE |

### Parámetros de Consulta para Listado de Usuarios

- `page`: Número de página (default: 1)
- `limit`: Cantidad de elementos por página (default: 100)
- `search`: Búsqueda por nombre o email

**Ejemplo**: `GET /api/v1/users?page=1&limit=10&search=juan`

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
PORT=:3000
RSA_PRIVATE_KEY_PATH=../api/certificates/app.rsa
RSA_PUBLIC_KEY_PATH=../api/certificates/app.rsa.pub
POSTGRES_USERNAME=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=appfedb
POSTGRES_DATABASE_URL=postgres://postgres:postgres@localhost:5432/appfedb?sslmode=disable
```

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

### Sin Docker

1. **Instalar PostgreSQL 16.2+**

2. **Configurar base de datos**
   ```sql
   CREATE DATABASE appfedb;
   CREATE USER postgres WITH PASSWORD 'postgres';
   GRANT ALL PRIVILEGES ON DATABASE appfedb TO postgres;
   ```

3. **Seguir pasos 1-3 y 5 de la instalación con Docker**

## 🔧 Comandos Útiles

```bash
# Instalar dependencias
go mod download

# Verificar y limpiar dependencias
go mod tidy

# Verificar integridad de módulos
go mod verify

# Compilar la aplicación
go build -o bin/api cmd/api/main.go

# Ejecutar tests
go test ./...

# Ejecutar con live reload (requiere air)
air
```

## 📊 Base de Datos

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

## 🛡️ Seguridad

- **Autenticación JWT** con algoritmo RSA256
- **Hash de contraseñas** con BCrypt (costo 12)
- **Validación de entrada** en todos los endpoints
- **Control de acceso basado en roles**
- **CORS configurado** para requests cross-origin
- **Middleware de seguridad** habilitado (Secure Headers, Gzip, etc.)

## 📝 Validaciones

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

1. **Inicialización**: Carga de variables de entorno y certificados RSA
2. **Conexión a BD**: Establecimiento de conexión con PostgreSQL
3. **Migraciones**: Ejecución automática de migraciones de base de datos
4. **Servicios**: Inicialización de servicios (User, Auth, JWT)
5. **Router**: Configuración de rutas y middleware
6. **Servidor**: Inicio del servidor HTTP con graceful shutdown

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto es parte de APPFE Lima y está sujeto a sus políticas internas de desarrollo.

## 🔗 Enlaces Relacionados

- **APPFE Lima**: [Sitio Web Oficial](https://appfelima.org)
- **Documentación de Echo**: [https://echo.labstack.com](https://echo.labstack.com)
- **PostgreSQL**: [https://postgresql.org](https://postgresql.org)

---

**Desarrollado con ❤️ para APPFE Lima**
