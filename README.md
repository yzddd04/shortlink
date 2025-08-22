# Link Shortener Backend

Sistem backend CRUD link shortener yang lengkap dengan Go, PostgreSQL, dan JWT authentication. Dibangun dengan arsitektur clean architecture dan best practices untuk scalability dan maintainability.

## Documentation Project

## Fitur Utama

- **CRUD Operations** - Create, Read, Update, Delete untuk link shortener
- **User Authentication** - JWT-based authentication dengan register/login
- **PostgreSQL Database** - Database yang robust dengan indexing
- **RESTful API** - API yang konsisten dan well-documented
- **Input Validation** - Validasi input yang comprehensive
- **Error Handling** - Error handling yang proper dan informatif
- **Rate Limiting** - Rate limiting untuk mencegah abuse
- **CORS Support** - Cross-origin resource sharing
- **Docker Support** - Containerization dengan Docker
- **Unit Testing** - Test coverage yang comprehensive
- **Security Features** - Password hashing, JWT tokens, input sanitization
- **URL Shortening** - Custom alias atau auto-generated short codes
- **Link Analytics** - Click tracking dan statistics
- **Link Expiration** - Expiration date untuk links
- **Pagination** - Pagination untuk list endpoints

## Arsitektur

Project ini menggunakan **Clean Architecture** dengan layer separation yang jelas:

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Handlers  │  │ Middleware  │  │      Routes         │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                     Business Layer                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Services  │  │   Models    │  │     Utils           │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                     Data Layer                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ Repository  │  │  Database   │  │    Migrations       │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Struktur Project

```
link-shortener/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── database.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── link.go
│   │   └── user.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── rate_limit.go
│   ├── models/
│   │   ├── link.go
│   │   └── user.go
│   ├── repository/
│   │   ├── link.go
│   │   └── user.go
│   ├── services/
│   │   ├── auth.go
│   │   ├── link.go
│   │   └── user.go
│   └── utils/
│       ├── jwt.go
│       └── validator.go
├── migrations/
│   └── 001_init.sql
├── docs/
│   └── api.md
├── tests/
│   ├── auth_test.go
│   └── link_test.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── env.example
└── README.md
```

## Quick Start

### Prerequisites

- **Go 1.21+** - Programming language
- **PostgreSQL 12+** - Database
- **Docker** (optional) - Containerization
- **Git** - Version control

### Installation

1. **Clone repository**
```bash
git clone <repository-url>
cd link-shortener
```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Setup environment**
```bash
cp env.example .env
# Edit .env file sesuai konfigurasi database
```

4. **Setup database**
```bash
# Jalankan PostgreSQL dan buat database
createdb link_shortener

# Jalankan migration
psql -d link_shortener -f migrations/001_init.sql
```

5. **Run application**
```bash
# Menggunakan Makefile
make run

# Atau langsung
go run cmd/server/main.go

# Atau menggunakan script
chmod +x scripts/run.sh
./scripts/run.sh
```

### Docker

```bash
# Build dan jalankan dengan Docker Compose
docker-compose up --build

# Atau build image terpisah
docker build -t link-shortener .
docker run -p 8080:8080 link-shortener
```

## API Documentation

### Authentication

#### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "user123",
  "email": "user@example.com",
  "password": "password123"
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Links

#### Create Short Link
```http
POST /api/links
Authorization: Bearer <token>
Content-Type: application/json

{
  "original_url": "https://example.com/very-long-url",
  "custom_alias": "my-link" // optional
}
```

#### Get All Links
```http
GET /api/links
Authorization: Bearer <token>
```

#### Get Link by ID
```http
GET /api/links/:id
Authorization: Bearer <token>
```

#### Update Link
```http
PUT /api/links/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "original_url": "https://new-url.com",
  "custom_alias": "new-alias"
}
```

#### Delete Link
```http
DELETE /api/links/:id
Authorization: Bearer <token>
```

#### Redirect to Original URL
```http
GET /r/:short_code
```

## Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test ./tests/auth_test.go

# Run with coverage
go test -cover ./...
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | password |
| DB_NAME | Database name | link_shortener |
| PORT | Server port | 8080 |
| JWT_SECRET | JWT secret key | - |
| JWT_EXPIRY | JWT expiry time | 24h |

## Contributing

1. Fork repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## License

MIT License
