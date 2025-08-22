# 🚀 Getting Started Guide

## 📋 Overview

Panduan lengkap untuk memulai development dengan Link Shortener Backend. Dokumen ini akan membantu Anda setup project dari awal hingga dapat menjalankan aplikasi.

## 🎯 What You'll Learn

Setelah menyelesaikan panduan ini, Anda akan dapat:
- ✅ Setup development environment
- ✅ Menjalankan aplikasi secara lokal
- ✅ Menggunakan Docker untuk development
- ✅ Testing API endpoints
- ✅ Memahami struktur project

## 📋 Prerequisites

Sebelum memulai, pastikan Anda telah menginstall:

### Required Software

1. **Go 1.21+**
   ```bash
   # Download dari https://golang.org/dl/
   # Atau menggunakan package manager
   
   # macOS (Homebrew)
   brew install go
   
   # Ubuntu/Debian
   sudo apt update
   sudo apt install golang-go
   
   # Windows
   # Download installer dari https://golang.org/dl/
   ```

2. **PostgreSQL 12+**
   ```bash
   # macOS (Homebrew)
   brew install postgresql
   brew services start postgresql
   
   # Ubuntu/Debian
   sudo apt update
   sudo apt install postgresql postgresql-contrib
   sudo systemctl start postgresql
   sudo systemctl enable postgresql
   
   # Windows
   # Download dari https://www.postgresql.org/download/windows/
   ```

3. **Git**
   ```bash
   # macOS (Homebrew)
   brew install git
   
   # Ubuntu/Debian
   sudo apt install git
   
   # Windows
   # Download dari https://git-scm.com/download/win
   ```

### Optional Software

4. **Docker & Docker Compose**
   ```bash
   # Download dari https://www.docker.com/products/docker-desktop
   ```

5. **VS Code dengan Extensions**
   - Go
   - REST Client
   - Docker
   - PostgreSQL

## 🛠️ Installation Steps

### Step 1: Clone Repository

```bash
# Clone repository
git clone <repository-url>
cd link-shortener

# Verify Go installation
go version
```

### Step 2: Setup Environment

```bash
# Install dependencies
go mod tidy
go mod download

# Copy environment file
cp env.example .env

# Edit environment variables
nano .env  # atau editor favorit Anda
```

**Contoh isi file .env:**
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=link_shortener
DB_SSL_MODE=disable

# Server Configuration
PORT=8080
GIN_MODE=debug

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRY=24h
```

### Step 3: Setup Database

```bash
# Create database
createdb link_shortener

# Atau menggunakan psql
psql -U postgres -c "CREATE DATABASE link_shortener;"

# Run migrations
psql -d link_shortener -f migrations/001_init.sql
```

### Step 4: Run Application

#### Option A: Using Scripts (Recommended)

**Linux/macOS:**
```bash
# Make script executable
chmod +x run.sh

# Run application
./run.sh
```

**Windows:**
```cmd
# Run application
run.bat
```

#### Option B: Using Makefile

```bash
# Install dependencies and run
make run

# Atau step by step
make deps
make build
make run
```

#### Option C: Direct Go Command

```bash
# Run directly
go run cmd/server/main.go

# Atau build first
go build -o bin/server cmd/server/main.go
./bin/server
```

### Step 5: Verify Installation

1. **Check Health Endpoint**
   ```bash
   curl http://localhost:8080/health
   ```

2. **Expected Response:**
   ```json
   {
     "status": "ok",
     "message": "Link Shortener API is running",
     "time": "2024-01-01T12:00:00Z"
   }
   ```

## 🐳 Docker Setup (Alternative)

Jika Anda lebih suka menggunakan Docker:

### Quick Start with Docker

```bash
# Build and run everything
docker-compose up --build

# Run in background
docker-compose up -d --build

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

### Access Services

- **Application**: http://localhost:8080
- **PgAdmin**: http://localhost:5050 (admin@example.com / admin)
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## 🧪 Testing the API

### 1. Register User

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. Create Link

```bash
# Replace YOUR_TOKEN with token from login response
curl -X POST http://localhost:8080/api/links \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "original_url": "https://example.com/very-long-url",
    "custom_alias": "test-link",
    "title": "Test Link"
  }'
```

### 4. Test Redirect

```bash
# Test the shortened link
curl -I http://localhost:8080/r/test-link
```

## 📚 Using the API Documentation

### REST Client (VS Code)

1. Install REST Client extension
2. Open `examples/api_test.http`
3. Replace `{{auth_token}}` with your JWT token
4. Click "Send Request" on any endpoint

### Postman

1. Import the collection from `examples/`
2. Set environment variables
3. Run requests

## 🔧 Development Workflow

### 1. Make Changes

```bash
# Create feature branch
git checkout -b feature/new-feature

# Make your changes
# Write tests
# Update documentation
```

### 2. Test Changes

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run linting
golangci-lint run
```

### 3. Commit Changes

```bash
# Add changes
git add .

# Commit with conventional commit message
git commit -m "feat: add new feature"

# Push to remote
git push origin feature/new-feature
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

3. **Permission Denied**
   ```bash
   # Fix file permissions
   chmod +x run.sh
   chmod +x scripts/run.sh
   ```

4. **Go Module Issues**
   ```bash
   # Clean module cache
   go clean -modcache
   
   # Re-download dependencies
   go mod download
   ```

### Getting Help

- Check the logs for error messages
- Verify all prerequisites are installed
- Ensure database is running and accessible
- Check environment variables are correct

## 📖 Next Steps

Setelah berhasil menjalankan aplikasi, Anda dapat:

1. **Explore the Codebase**
   - Read through the handlers
   - Understand the service layer
   - Check the repository pattern

2. **Add New Features**
   - Implement new endpoints
   - Add validation rules
   - Create new models

3. **Improve the Application**
   - Add more tests
   - Implement caching
   - Add monitoring

4. **Deploy to Production**
   - Read deployment guide
   - Set up CI/CD
   - Configure production environment

## 📚 Additional Resources

- [API Documentation](api.md)
- [Development Guide](development.md)
- [Deployment Guide](deployment.md)
- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)

## 🤝 Need Help?

Jika Anda mengalami masalah:

1. Check the troubleshooting section above
2. Review the logs for error messages
3. Verify all prerequisites are installed
4. Check the documentation
5. Create an issue on GitHub

---

**Selamat!** 🎉 Anda telah berhasil setup dan menjalankan Link Shortener Backend. Sekarang Anda siap untuk development!
