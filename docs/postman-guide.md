# 📮 Panduan Lengkap Postman untuk Link Shortener API

## 📋 Overview

Panduan ini akan membantu Anda menggunakan Postman untuk testing Link Shortener API dengan langkah-langkah yang detail dan mudah diikuti.

## 🚀 Setup Awal

### 1. **Prerequisites**

Sebelum menggunakan Postman, pastikan:

- ✅ **Aplikasi sudah berjalan** di `http://localhost:8080`
- ✅ **PostgreSQL sudah running** dan database sudah dibuat
- ✅ **Postman sudah terinstall** di komputer Anda

### 2. **Jalankan Aplikasi**

```bash
# Option 1: Menggunakan script
./run.sh  # Linux/macOS
run.bat   # Windows

# Option 2: Menggunakan Makefile
make run

# Option 3: Direct Go command
go run cmd/server/main.go
```

### 3. **Verifikasi Aplikasi Berjalan**

Buka browser dan akses: `http://localhost:8080/health`

Anda harus melihat response:
```json
{
  "status": "ok",
  "message": "Link Shortener API is running",
  "time": "2024-01-01T12:00:00Z"
}
```

## 📥 Import Collection dan Environment

### 1. **Download Files**

Download file-file berikut dari project:
- `examples/Link_Shortener_API.postman_collection.json`
- `examples/Link_Shortener_Environment.postman_environment.json`

### 2. **Import ke Postman**

#### **Import Collection:**
1. Buka Postman
2. Klik **"Import"** button
3. Drag & drop file `Link_Shortener_API.postman_collection.json`
4. Klik **"Import"**

#### **Import Environment:**
1. Klik **"Import"** button lagi
2. Drag & drop file `Link_Shortener_Environment.postman_environment.json`
3. Klik **"Import"**

### 3. **Setup Environment**

1. Di dropdown environment (kanan atas), pilih **"Link Shortener Environment"**
2. Klik **"View more actions"** (titik tiga) → **"Edit"**
3. Pastikan variable `base_url` = `http://localhost:8080`
4. Klik **"Save"**

## 🧪 Testing Workflow

### **Step 1: Health Check**

1. Buka collection **"Link Shortener API"**
2. Klik **"Health Check"**
3. Klik **"Send"**

**Expected Response:**
```json
{
  "status": "ok",
  "message": "Link Shortener API is running",
  "time": "2024-01-01T12:00:00Z"
}
```

### **Step 2: Register User**

1. Buka folder **"Authentication"**
2. Klik **"Register User"**
3. Klik **"Send"**

**Expected Response:**
```json
{
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "uuid",
      "username": "testuser",
      "email": "test@example.com",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    },
    "token": "jwt-token-here"
  }
}
```

### **Step 3: Login User**

1. Klik **"Login User"**
2. Klik **"Send"**

**Expected Response:**
```json
{
  "message": "Login successful",
  "data": {
    "user": {
      "id": "uuid",
      "username": "testuser",
      "email": "test@example.com",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    },
    "token": "jwt-token-here"
  }
}
```

**🎯 Important:** Token akan otomatis tersimpan di environment variable `auth_token`

### **Step 4: Create Link**

1. Buka folder **"Links"**
2. Klik **"Create Link"**
3. Klik **"Send"**

**Expected Response:**
```json
{
  "message": "Link created successfully",
  "data": {
    "id": "uuid",
    "original_url": "https://example.com/very-long-url",
    "short_code": "test-link",
    "short_url": "http://localhost:8080/r/test-link",
    "title": "Test Link",
    "clicks": 0,
    "is_active": true,
    "expires_at": "2024-12-31T23:59:59Z",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

**🎯 Important:** Copy `id` dari response untuk digunakan di request berikutnya

### **Step 5: Update Environment Variables**

1. Klik **"View more actions"** (titik tiga) → **"Edit"**
2. Set `link_id` = ID dari response Create Link
3. Set `short_code` = short_code dari response (default: "test-link")
4. Klik **"Save"**

### **Step 6: Test Other Endpoints**

Sekarang Anda bisa test endpoint lainnya:

#### **Get All Links**
- Klik **"Get All Links"** → **"Send"**

#### **Get Link by ID**
- Klik **"Get Link by ID"** → **"Send"**

#### **Update Link**
- Klik **"Update Link"** → **"Send"**

#### **Get Link Statistics**
- Klik **"Get Link Statistics"** → **"Send"**

#### **Test Redirect**
- Klik **"Redirect to Original URL"** → **"Send"**

#### **Delete Link**
- Klik **"Delete Link"** → **"Send"**

## 🔧 Manual Setup (Tanpa Import)

Jika Anda ingin setup manual tanpa import file:

### 1. **Create New Collection**

1. Klik **"New"** → **"Collection"**
2. Name: `Link Shortener API`
3. Klik **"Create"**

### 2. **Create Environment**

1. Klik **"New"** → **"Environment"**
2. Name: `Link Shortener Environment`
3. Add variables:
   - `base_url`: `http://localhost:8080`
   - `auth_token`: (kosong)
   - `link_id`: (kosong)
   - `short_code`: `test-link`
4. Klik **"Save"**

### 3. **Create Requests**

#### **Health Check**
```
Method: GET
URL: {{base_url}}/health
```

#### **Register User**
```
Method: POST
URL: {{base_url}}/api/auth/register
Headers: Content-Type: application/json
Body (raw JSON):
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

#### **Login User**
```
Method: POST
URL: {{base_url}}/api/auth/login
Headers: Content-Type: application/json
Body (raw JSON):
{
  "email": "test@example.com",
  "password": "password123"
}
```

#### **Create Link**
```
Method: POST
URL: {{base_url}}/api/links
Headers: 
  Content-Type: application/json
  Authorization: Bearer {{auth_token}}
Body (raw JSON):
{
  "original_url": "https://example.com/very-long-url",
  "custom_alias": "test-link",
  "title": "Test Link"
}
```

## 🎯 Tips dan Best Practices

### 1. **Token Management**
- Token akan otomatis tersimpan setelah login
- Jika token expired, login ulang untuk mendapatkan token baru

### 2. **Environment Variables**
- Gunakan `{{variable_name}}` untuk menggunakan environment variables
- Update `link_id` setelah create link untuk testing endpoint lainnya

### 3. **Error Handling**
- Test folder **"Error Testing"** untuk melihat error responses
- Perhatikan status code dan error messages

### 4. **Request Headers**
- Semua request yang memerlukan auth harus include: `Authorization: Bearer {{auth_token}}`
- POST/PUT requests harus include: `Content-Type: application/json`

### 5. **Response Validation**
- Check status code (200, 201, 400, 401, 404, etc.)
- Validate response structure
- Check error messages untuk debugging

## 🐛 Troubleshooting

### **Common Issues:**

#### 1. **Connection Refused**
```
Error: connect ECONNREFUSED 127.0.0.1:8080
```
**Solution:** Pastikan aplikasi berjalan di port 8080

#### 2. **401 Unauthorized**
```
Error: Authorization header required
```
**Solution:** 
- Login ulang untuk mendapatkan token baru
- Pastikan header `Authorization: Bearer {{auth_token}}` ada

#### 3. **400 Bad Request**
```
Error: Invalid request data
```
**Solution:** 
- Check JSON format di request body
- Pastikan semua required fields terisi

#### 4. **404 Not Found**
```
Error: Link not found
```
**Solution:** 
- Update `link_id` variable dengan ID yang valid
- Pastikan link belum di-delete

### **Debug Steps:**

1. **Check Application Status**
   - Test Health Check endpoint
   - Check application logs

2. **Check Environment Variables**
   - Verify `base_url` = `http://localhost:8080`
   - Check `auth_token` tidak kosong

3. **Check Request Format**
   - Verify HTTP method
   - Check headers
   - Validate JSON body

4. **Check Database**
   - Pastikan PostgreSQL running
   - Verify database connection

## 📊 Testing Scenarios

### **Complete Workflow Test:**

1. **Health Check** → Should return 200 OK
2. **Register User** → Should return 201 Created
3. **Login User** → Should return 200 OK + token
4. **Create Link** → Should return 201 Created
5. **Get All Links** → Should return 200 OK + links array
6. **Get Link by ID** → Should return 200 OK + link details
7. **Update Link** → Should return 200 OK + updated link
8. **Get Statistics** → Should return 200 OK + stats
9. **Test Redirect** → Should return 301/302 redirect
10. **Delete Link** → Should return 200 OK

### **Error Scenarios:**

1. **Invalid Authentication** → Should return 401
2. **Missing Authentication** → Should return 401
3. **Invalid Link ID** → Should return 400
4. **Duplicate Email** → Should return 400
5. **Invalid JSON** → Should return 400

## 🎉 Success Criteria

Anda berhasil menggunakan Postman jika:

- ✅ Semua endpoints berhasil di-test
- ✅ Response sesuai dengan expected format
- ✅ Error handling berfungsi dengan baik
- ✅ Environment variables bekerja dengan benar
- ✅ Token management berfungsi otomatis

---

**Selamat!** 🎉 Anda telah berhasil menggunakan Postman untuk testing Link Shortener API. Sekarang Anda siap untuk development dan testing yang lebih advanced!
