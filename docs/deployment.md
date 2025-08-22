# Deployment Guide

## Overview

This document provides comprehensive instructions for deploying the Link Shortener Backend in various environments.

## Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Docker (optional)
- Git

## Local Development Setup

### 1. Clone Repository
```bash
git clone <repository-url>
cd link-shortener
```

### 2. Install Dependencies
```bash
go mod tidy
go mod download
```

### 3. Environment Configuration
```bash
cp env.example .env
# Edit .env file with your configuration
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
# Using Makefile
make run

# Or directly
go run cmd/server/main.go

# Or using script
chmod +x scripts/run.sh
./scripts/run.sh
```

## Docker Deployment

### 1. Build Image
```bash
docker build -t link-shortener .
```

### 2. Run with Docker Compose
```bash
docker-compose up --build
```

### 3. Environment Variables for Docker
Create a `.env` file for Docker:
```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=link_shortener
DB_SSL_MODE=disable
PORT=8080
GIN_MODE=release
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRY=24h
```

## Production Deployment

### 1. Environment Configuration
```bash
# Production environment variables
export GIN_MODE=release
export DB_HOST=your-db-host
export DB_PORT=5432
export DB_USER=your-db-user
export DB_PASSWORD=your-secure-password
export DB_NAME=link_shortener
export DB_SSL_MODE=require
export JWT_SECRET=your-very-secure-jwt-secret
export JWT_EXPIRY=24h
```

### 2. Build for Production
```bash
# Build binary
GOOS=linux GOARCH=amd64 go build -o bin/server cmd/server/main.go

# Or build with optimizations
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/server cmd/server/main.go
```

### 3. Database Setup
```sql
-- Create production database
CREATE DATABASE link_shortener;

-- Create user with limited privileges
CREATE USER link_shortener_user WITH PASSWORD 'secure_password';
GRANT CONNECT ON DATABASE link_shortener TO link_shortener_user;
GRANT USAGE ON SCHEMA public TO link_shortener_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO link_shortener_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO link_shortener_user;
```

### 4. Run with Process Manager
```bash
# Using systemd
sudo cp systemd/link-shortener.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable link-shortener
sudo systemctl start link-shortener

# Using PM2
npm install -g pm2
pm2 start ecosystem.config.js
```

## Cloud Deployment

### AWS EC2
```bash
# Install Go
wget https://golang.org/dl/go1.21.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# Setup application
git clone <repository-url>
cd link-shortener
make setup
make build
```

### Google Cloud Run
```bash
# Build and push to Container Registry
gcloud builds submit --tag gcr.io/PROJECT_ID/link-shortener

# Deploy to Cloud Run
gcloud run deploy link-shortener \
  --image gcr.io/PROJECT_ID/link-shortener \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars "DB_HOST=your-db-host,DB_PORT=5432"
```

### Heroku
```bash
# Create Heroku app
heroku create your-link-shortener-app

# Add PostgreSQL addon
heroku addons:create heroku-postgresql:hobby-dev

# Set environment variables
heroku config:set GIN_MODE=release
heroku config:set JWT_SECRET=your-secret-key

# Deploy
git push heroku main
```

## Monitoring and Logging

### 1. Health Checks
```bash
# Check application health
curl http://localhost:8080/health
```

### 2. Logging
```bash
# View application logs
tail -f /var/log/link-shortener/app.log

# Using systemd
journalctl -u link-shortener -f
```

### 3. Metrics
```bash
# Monitor database connections
psql -d link_shortener -c "SELECT count(*) FROM pg_stat_activity;"

# Monitor application performance
curl http://localhost:8080/api/links/stats
```

## Security Considerations

### 1. Environment Variables
- Never commit `.env` files to version control
- Use strong, unique passwords
- Rotate JWT secrets regularly

### 2. Database Security
- Use SSL connections in production
- Implement connection pooling
- Regular backups

### 3. Network Security
- Use HTTPS in production
- Implement rate limiting
- Configure CORS properly

### 4. Application Security
- Keep dependencies updated
- Implement input validation
- Use prepared statements for SQL

## Troubleshooting

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
   chmod +x bin/server
   chmod +x scripts/run.sh
   ```

### Log Analysis
```bash
# View error logs
grep ERROR /var/log/link-shortener/app.log

# Monitor real-time logs
tail -f /var/log/link-shortener/app.log | grep -E "(ERROR|WARN)"
```

## Performance Optimization

### 1. Database Optimization
```sql
-- Create indexes for better performance
CREATE INDEX CONCURRENTLY idx_links_user_id_created_at ON links(user_id, created_at DESC);
CREATE INDEX CONCURRENTLY idx_links_short_code_active ON links(short_code) WHERE is_active = true;
```

### 2. Application Optimization
```bash
# Build with optimizations
go build -ldflags="-s -w" -o bin/server cmd/server/main.go

# Use connection pooling
export DB_MAX_OPEN_CONNS=25
export DB_MAX_IDLE_CONNS=5
```

### 3. Caching
```bash
# Install Redis for caching
sudo apt install redis-server

# Configure Redis in application
export REDIS_HOST=localhost
export REDIS_PORT=6379
```

## Backup and Recovery

### 1. Database Backup
```bash
# Create backup
pg_dump link_shortener > backup_$(date +%Y%m%d_%H%M%S).sql

# Restore backup
psql link_shortener < backup_20240101_120000.sql
```

### 2. Application Backup
```bash
# Backup configuration
tar -czf config_backup_$(date +%Y%m%d).tar.gz .env migrations/

# Backup logs
tar -czf logs_backup_$(date +%Y%m%d).tar.gz /var/log/link-shortener/
```

## Scaling

### 1. Horizontal Scaling
```bash
# Load balancer configuration
upstream link_shortener {
    server 127.0.0.1:8080;
    server 127.0.0.1:8081;
    server 127.0.0.1:8082;
}
```

### 2. Database Scaling
```bash
# Read replicas
export DB_READ_HOST=read-replica-host
export DB_WRITE_HOST=master-host
```

### 3. Caching Layer
```bash
# Redis cluster
export REDIS_CLUSTER=true
export REDIS_NODES="node1:6379,node2:6379,node3:6379"
```

