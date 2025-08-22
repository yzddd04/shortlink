# 🛠️ Development Guide

## 📋 Overview

Dokumen ini berisi panduan lengkap untuk development Link Shortener Backend, termasuk setup environment, best practices, dan workflow development.

## 🚀 Quick Start Development

### 1. Prerequisites

Pastikan Anda telah menginstall:
- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **PostgreSQL 12+** - [Download PostgreSQL](https://www.postgresql.org/download/)
- **Git** - [Download Git](https://git-scm.com/downloads)
- **Docker** (optional) - [Download Docker](https://www.docker.com/products/docker-desktop)

### 2. Clone Repository

```bash
git clone <repository-url>
cd link-shortener
```

### 3. Setup Development Environment

```bash
# Install dependencies
go mod tidy
go mod download

# Copy environment file
cp env.example .env

# Edit environment variables
nano .env  # atau editor favorit Anda
```

### 4. Database Setup

```bash
# Create database
createdb link_shortener

# Run migrations
psql -d link_shortener -f migrations/001_init.sql
```

### 5. Run Application

```bash
# Development mode
go run cmd/server/main.go

# Atau menggunakan Makefile
make run

# Atau menggunakan script
./scripts/run.sh  # Linux/macOS
scripts\run.bat   # Windows
```

## 🔧 Development Tools

### Recommended IDE/Editor

- **VS Code** dengan extensions:
  - Go
  - REST Client
  - Docker
  - PostgreSQL

- **GoLand** (JetBrains)
- **Vim/Neovim** dengan Go plugins

### Useful Development Tools

```bash
# Install development tools
go install github.com/cosmtrek/air@latest          # Hot reload
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest  # Linter
go install github.com/go-delve/delve/cmd/dlv@latest  # Debugger
go install github.com/vektra/mockery/v2@latest     # Mock generator
```

## 📁 Project Structure

```
link-shortener/
├── cmd/                    # Application entry points
│   └── server/
│       └── main.go        # Main application
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── database/         # Database connection
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Data models
│   ├── repository/       # Data access layer
│   ├── services/         # Business logic
│   └── utils/            # Utility functions
├── migrations/           # Database migrations
├── docs/                 # Documentation
├── tests/                # Test files
├── scripts/              # Build and deployment scripts
├── examples/             # Example files
└── docker-compose.yml    # Docker configuration
```

## 🧪 Testing

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./tests/auth_test.go

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...
```

### Test Structure

```go
// Example test structure
func TestCreateLink(t *testing.T) {
    // Arrange
    // Setup test data and dependencies
    
    // Act
    // Execute the function being tested
    
    // Assert
    // Verify the results
}
```

### Mocking

```bash
# Generate mocks
mockery --dir internal/repository --name UserRepository
mockery --dir internal/services --name AuthService
```

## 🔍 Debugging

### Using Delve Debugger

```bash
# Start debugger
dlv debug cmd/server/main.go

# Common commands
break main.main          # Set breakpoint
continue                 # Continue execution
next                     # Step over
step                     # Step into
print variable           # Print variable value
```

### Using VS Code

1. Install Go extension
2. Set breakpoints in code
3. Press F5 to start debugging
4. Use debug console for evaluation

### Logging

```go
// Add logging to your code
log.Printf("Processing request: %s", requestID)
log.Printf("Database query executed: %s", query)
log.Printf("Error occurred: %v", err)
```

## 📝 Code Style

### Go Formatting

```bash
# Format code
go fmt ./...

# Or using Makefile
make fmt
```

### Linting

```bash
# Run linter
golangci-lint run

# Or using Makefile
make lint
```

### Code Review Checklist

- [ ] Code follows Go conventions
- [ ] Functions are properly documented
- [ ] Error handling is implemented
- [ ] Tests are written and passing
- [ ] No hardcoded values
- [ ] Proper logging is added
- [ ] Security considerations are addressed

## 🔄 Development Workflow

### 1. Feature Development

```bash
# Create feature branch
git checkout -b feature/new-feature

# Make changes
# Write tests
# Update documentation

# Commit changes
git add .
git commit -m "feat: add new feature"

# Push to remote
git push origin feature/new-feature
```

### 2. Code Review

- Create Pull Request
- Request review from team members
- Address feedback
- Merge after approval

### 3. Testing

```bash
# Run tests before commit
make test

# Run linting
make lint

# Check formatting
make fmt
```

## 🐳 Docker Development

### Local Development with Docker

```bash
# Start services
docker-compose up -d postgres redis

# Run application locally
go run cmd/server/main.go
```

### Full Docker Development

```bash
# Build and run everything
docker-compose up --build

# View logs
docker-compose logs -f app

# Access services
# App: http://localhost:8080
# PgAdmin: http://localhost:5050
# Redis: localhost:6379
```

## 🔧 Configuration

### Environment Variables

```bash
# Development
GIN_MODE=debug
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=link_shortener
JWT_SECRET=dev-secret-key

# Production
GIN_MODE=release
DB_SSL_MODE=require
JWT_SECRET=very-secure-production-key
```

### Database Configuration

```sql
-- Development database setup
CREATE DATABASE link_shortener;
CREATE USER link_dev WITH PASSWORD 'dev_password';
GRANT ALL PRIVILEGES ON DATABASE link_shortener TO link_dev;
```

## 🚀 Hot Reload Development

### Using Air

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Create .air.toml configuration
air init

# Start with hot reload
air
```

### Air Configuration (.air.toml)

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false
```

## 📊 Performance Monitoring

### Profiling

```bash
# CPU profiling
go test -cpuprofile cpu.prof ./...

# Memory profiling
go test -memprofile mem.prof ./...

# Analyze profiles
go tool pprof cpu.prof
go tool pprof mem.prof
```

### Benchmarking

```bash
# Run benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkCreateLink ./tests/
```

## 🔒 Security Development

### Security Best Practices

1. **Input Validation**
   ```go
   // Always validate input
   if err := validateURL(req.OriginalURL); err != nil {
       return err
   }
   ```

2. **SQL Injection Prevention**
   ```go
   // Use prepared statements
   query := `SELECT * FROM users WHERE id = $1`
   row := db.QueryRow(query, userID)
   ```

3. **Password Hashing**
   ```go
   // Use bcrypt for password hashing
   hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
   ```

4. **JWT Security**
   ```go
   // Use strong secret keys
   // Set appropriate expiration times
   // Validate token claims
   ```

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   ```bash
   # Check PostgreSQL status
   sudo systemctl status postgresql
   
   # Check connection
   psql -h localhost -U postgres -d link_shortener
   ```

2. **Port Already in Use**
   ```bash
   # Find process using port
   lsof -i :8080
   
   # Kill process
   kill -9 <PID>
   ```

3. **Go Module Issues**
   ```bash
   # Clean module cache
   go clean -modcache
   
   # Re-download dependencies
   go mod download
   ```

4. **Docker Issues**
   ```bash
   # Clean Docker
   docker system prune -a
   
   # Rebuild images
   docker-compose build --no-cache
   ```

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [JWT Documentation](https://jwt.io/)
- [Docker Documentation](https://docs.docker.com/)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write tests
5. Update documentation
6. Submit a pull request

### Commit Message Format

```
type(scope): description

feat(auth): add JWT authentication
fix(api): resolve CORS issue
docs(readme): update installation guide
test(handlers): add unit tests for link creation
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build process changes
