# GopherSocial

API REST de una red social construida en Go. Los usuarios pueden registrarse, crear posts, seguir a otros, comentar y recibir un feed personalizado.

**API en producción:** [gophersocial.onrender.com](https://gophersocial.onrender.com/v1/health)  
**Frontend:** [gopher-social-web.vercel.app](https://gopher-social-web.vercel.app)  
**Documentación:** [gophersocial.onrender.com/v1/swagger](https://gophersocial.onrender.com/v1/swagger/index.html)

---

## Funcionalidades

- Autenticación con JWT y flujo de registro por email
- Confirmación de cuenta vía email transaccional (Brevo)
- Posts con comentarios y feed personalizado según usuarios seguidos
- Sistema de seguir / dejar de seguir con usuarios sugeridos
- Control de acceso por roles (user / moderador / admin)
- Capa de caché con Redis para consultas de usuarios
- Middleware de rate limiting
- Documentación Swagger
- CI con GitHub Actions (build, vet, staticcheck, tests)
- CD con Render (deploy automático al hacer push a `main`)

## Stack Tecnológico

| Capa | Tecnología |
|------|------------|
| Lenguaje | [Go](https://go.dev/) |
| Router | [chi](https://github.com/go-chi/chi) |
| Base de datos | [PostgreSQL 16](https://www.postgresql.org/) |
| Caché | [Redis](https://redis.io/) via [go-redis](https://github.com/redis/go-redis) |
| Autenticación | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) |
| Migraciones | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Email | [Brevo](https://www.brevo.com/) |
| Documentación | [Swagger / swaggo](https://github.com/swaggo/swag) |
| Logging | [zap](https://github.com/uber-go/zap) |
| Validación | [go-playground/validator](https://github.com/go-playground/validator) |
| Contenedores | [Docker](https://www.docker.com/) |

## Estructura del Proyecto

```
.
├── cmd/
│   ├── api/          # Handlers HTTP, middleware, rutas
│   └── migrate/      # Migraciones y datos de prueba
├── internal/
│   ├── auth/         # Autenticador JWT
│   ├── db/           # Conexión a la base de datos
│   ├── env/          # Helpers para variables de entorno
│   ├── mailer/       # Cliente de email y templates
│   ├── ratelimiter/  # Rate limiting
│   └── store/        # Capa de acceso a datos (PostgreSQL + Redis)
└── docs/             # Documentación Swagger autogenerada
```

## Endpoints

| Método | Endpoint | Descripción | Auth |
|--------|----------|-------------|------|
| POST | `/v1/auth/register` | Registrar un nuevo usuario | No |
| POST | `/v1/auth/token` | Login y obtener JWT | No |
| PUT | `/v1/users/activate/{token}` | Activar cuenta por email | No |
| GET | `/v1/users/{id}` | Obtener perfil de usuario | JWT |
| GET | `/v1/users/search?q=` | Buscar usuarios por username | JWT |
| GET | `/v1/users/suggested` | Usuarios sugeridos para seguir | JWT |
| PUT | `/v1/users/{id}/follow` | Seguir a un usuario | JWT |
| PUT | `/v1/users/{id}/unfollow` | Dejar de seguir a un usuario | JWT |
| GET | `/v1/feed` | Obtener feed personalizado | JWT |
| POST | `/v1/posts` | Crear un post | JWT |
| GET | `/v1/posts/{id}` | Obtener un post | JWT |
| PATCH | `/v1/posts/{id}` | Actualizar un post | JWT (dueño/moderador) |
| DELETE | `/v1/posts/{id}` | Eliminar un post | JWT (dueño/admin) |
| POST | `/v1/posts/{id}/comments` | Agregar un comentario | JWT |

## Cómo correr el proyecto

### Requisitos

- Go 1.21+
- Docker
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### Correr localmente

```bash
# Levantar PostgreSQL y Redis
docker compose up -d

# Copiar y completar las variables de entorno
cp .envrc.example .envrc
source .envrc

# Correr las migraciones
make migrate-up

# Cargar datos de prueba (opcional)
make seed

# Iniciar el servidor
go run ./api
```

