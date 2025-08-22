@echo off
REM Link Shortener Backend Runner Script for Windows

echo 🚀 Starting Link Shortener Backend...

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go is not installed. Please install Go 1.21+ first.
    pause
    exit /b 1
)

echo ✅ Go is installed

REM Install dependencies
echo 📦 Installing dependencies...
go mod tidy
go mod download

REM Check if .env file exists
if not exist .env (
    echo 📝 Creating .env file from env.example...
    copy env.example .env
    echo ⚠️  Please edit .env file with your database configuration
)

REM Check if PostgreSQL is running (optional check)
echo ⚠️  Please ensure PostgreSQL is running

REM Build the application
echo 🔨 Building application...
go build -o bin\server.exe cmd\server\main.go

REM Run the application
echo 🎯 Starting server...
echo 📍 Server will be available at: http://localhost:8080
echo 📚 API Documentation: http://localhost:8080/health
echo.
echo Press Ctrl+C to stop the server
echo.

bin\server.exe

pause

