# 🔗 Link Shortener API Documentation

## 📋 Overview

Link Shortener API adalah RESTful API yang menyediakan layanan untuk membuat, mengelola, dan mengakses link pendek. API ini menggunakan JWT untuk autentikasi dan mendukung berbagai fitur seperti custom alias, analytics, dan link expiration.

## 🌐 Base URL

```
http://localhost:8080
```

## 🔐 Authentication

Sebagian besar endpoint memerlukan autentikasi JWT. Sertakan token dalam header Authorization:

```
Authorization: Bearer <your-jwt-token>
```

### Token Format
- **Type**: JWT (JSON Web Token)
- **Algorithm**: HS256
- **Expiration**: 24 jam (dapat dikonfigurasi)
- **Payload**: User ID, Username, Email

## Endpoints

### Health Check
**GET** `/health`

Check if the API is running.

**Response:**
```json
{
  "status": "ok",
  "message": "Link Shortener API is running",
  "time": "2024-01-01T12:00:00Z"
}
```

### Authentication

#### Register User
**POST** `/api/auth/register`

Register a new user account.

**Request Body:**
```json
{
  "username": "user123",
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "uuid",
      "username": "user123",
      "email": "user@example.com",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    },
    "token": "jwt-token"
  }
}
```

#### Login
**POST** `/api/auth/login`

Login with existing credentials.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "data": {
    "user": {
      "id": "uuid",
      "username": "user123",
      "email": "user@example.com",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    },
    "token": "jwt-token"
  }
}
```

#### Get Profile
**GET** `/api/auth/profile`

Get current user profile (requires authentication).

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "username": "user123",
    "email": "user@example.com",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

### Links

#### Create Link
**POST** `/api/links`

Create a new shortened link (requires authentication).

**Request Body:**
```json
{
  "original_url": "https://example.com/very-long-url",
  "custom_alias": "my-link",
  "title": "My Custom Link",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

**Response:**
```json
{
  "message": "Link created successfully",
  "data": {
    "id": "uuid",
    "original_url": "https://example.com/very-long-url",
    "short_code": "my-link",
    "short_url": "http://localhost:8080/r/my-link",
    "title": "My Custom Link",
    "clicks": 0,
    "is_active": true,
    "expires_at": "2024-12-31T23:59:59Z",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### Get All Links
**GET** `/api/links`

Get all links for the authenticated user with pagination.

**Query Parameters:**
- `limit` (optional): Number of links per page (default: 10, max: 100)
- `offset` (optional): Number of links to skip (default: 0)

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "original_url": "https://example.com/very-long-url",
      "short_code": "my-link",
      "short_url": "http://localhost:8080/r/my-link",
      "title": "My Custom Link",
      "clicks": 5,
      "is_active": true,
      "expires_at": null,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ],
  "pagination": {
    "limit": 10,
    "offset": 0
  }
}
```

#### Get Link by ID
**GET** `/api/links/:id`

Get a specific link by ID (requires authentication).

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "original_url": "https://example.com/very-long-url",
    "short_code": "my-link",
    "short_url": "http://localhost:8080/r/my-link",
    "title": "My Custom Link",
    "clicks": 5,
    "is_active": true,
    "expires_at": null,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### Update Link
**PUT** `/api/links/:id`

Update a specific link (requires authentication).

**Request Body:**
```json
{
  "original_url": "https://new-url.com",
  "custom_alias": "new-alias",
  "title": "Updated Link Title",
  "is_active": false,
  "expires_at": "2024-12-31T23:59:59Z"
}
```

**Response:**
```json
{
  "message": "Link updated successfully",
  "data": {
    "id": "uuid",
    "original_url": "https://new-url.com",
    "short_code": "new-alias",
    "short_url": "http://localhost:8080/r/new-alias",
    "title": "Updated Link Title",
    "clicks": 5,
    "is_active": false,
    "expires_at": "2024-12-31T23:59:59Z",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### Delete Link
**DELETE** `/api/links/:id`

Delete a specific link (requires authentication).

**Response:**
```json
{
  "message": "Link deleted successfully"
}
```

#### Get Link Statistics
**GET** `/api/links/stats`

Get statistics for the authenticated user's links.

**Response:**
```json
{
  "data": {
    "total_links": 10,
    "total_clicks": 150,
    "active_links": 8,
    "expired_links": 2
  }
}
```

### Redirect

#### Redirect to Original URL
**GET** `/r/:shortCode`

Redirect to the original URL using the short code (public endpoint).

**Response:** HTTP 301 redirect to the original URL.

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request data",
  "details": "Field validation error details"
}
```

### 401 Unauthorized
```json
{
  "error": "Authorization header required"
}
```

### 404 Not Found
```json
{
  "error": "Link not found or expired"
}
```

### 429 Too Many Requests
```json
{
  "error": "Rate limit exceeded"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

## Rate Limiting

The API implements rate limiting:
- 100 requests per minute per IP address
- Rate limit headers are included in responses

## Validation Rules

### User Registration
- Username: 3-50 characters, unique
- Email: Valid email format, unique
- Password: Minimum 6 characters

### Link Creation
- Original URL: Valid URL format
- Custom Alias: 3-20 characters, alphanumeric and hyphens only, unique
- Title: Maximum 255 characters
- Expires At: Valid future date (optional)

## Examples

### Complete Workflow

1. **Register a user:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

2. **Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

3. **Create a link:**
```bash
curl -X POST http://localhost:8080/api/links \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "original_url": "https://example.com/very-long-url",
    "custom_alias": "my-link"
  }'
```

4. **Access the shortened link:**
```bash
curl -I http://localhost:8080/r/my-link
```
