#!/bin/bash

# Link Shortener Backend Runner Script

set -e

echo "🚀 Starting Link Shortener Backend..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ Go version $GO_VERSION is too old. Please install Go $REQUIRED_VERSION+"
    exit 1
fi

echo "✅ Go version: $GO_VERSION"

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy
go mod download

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from env.example..."
    cp env.example .env
    echo "⚠️  Please edit .env file with your database configuration"
fi

# Check if PostgreSQL is running (optional check)
if command -v pg_isready &> /dev/null; then
    if pg_isready -q; then
        echo "✅ PostgreSQL is running"
    else
        echo "⚠️  PostgreSQL might not be running. Please ensure it's started."
    fi
else
    echo "⚠️  pg_isready not found. Please ensure PostgreSQL is running."
fi

# Build the application
echo "🔨 Building application..."
go build -o bin/server cmd/server/main.go

# Run the application
echo "🎯 Starting server..."
echo "📍 Server will be available at: http://localhost:8080"
echo "📚 API Documentation: http://localhost:8080/health"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

./bin/server

