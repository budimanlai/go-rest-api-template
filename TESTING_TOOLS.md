# Testing Tools Documentation

Tools untuk testing sistem JWT authentication yang telah diupdate dan siap digunakan.

## 1. Postman Collection - `postman_collection.json`

Collection lengkap dengan endpoint authentication terbaru:
- **GET** `/api/v1/public/auth/token` - Get public token
- **POST** `/api/v1/public/auth/register` - Register user  
- **POST** `/api/v1/public/auth/login` - Login user
- **GET** `/api/v1/private/users` - Access private endpoint
- **POST** `/api/v1/public/auth/refresh` - Refresh token
- **POST** `/api/v1/public/auth/logout` - Logout

Plus error testing scenarios.

Variables yang sudah dikonfigurasi:
- `base_url`: http://localhost:3000
- `api_key`: test-api-key-12345
- `public_token`: (akan diisi otomatis)
- `private_token`: (akan diisi otomatis)

## 2. Manual Test Guide - `manual_test_guide.sh`

Script shell yang dapat:
- Menampilkan semua command cURL manual
- Menjalankan automated testing dengan `./manual_test_guide.sh test`

Usage:
```bash
chmod +x manual_test_guide.sh

# Lihat manual commands
./manual_test_guide.sh

# Jalankan automated test
./manual_test_guide.sh test
```

## 3. JWT Tester CLI - `jwt_tester.go`

Program Go yang dapat dijalankan langsung:
```bash
go run jwt_tester.go
```

Testing flow lengkap:
1. Get public token dengan API key
2. Register user (skip jika sudah ada)
3. Login untuk mendapat private token
4. Test access ke private endpoint
5. Test refresh token
6. Test logout
7. Error testing (invalid API key, wrong token type, etc.)

## Authentication Flow

```
1. API Key Only → GET /api/v1/public/auth/token → Public JWT Token
2. Public JWT + API Key → POST /api/v1/public/auth/register → User Created
3. Public JWT + API Key → POST /api/v1/public/auth/login → Private JWT Token
4. Private JWT + API Key → GET /api/v1/users → Protected User Data (requires private token)
```

**Note**: All user endpoints now require Private JWT middleware (API Key + Private Token)

## Environment Setup

Pastikan API key `test-api-key-12345` sudah ada di database atau sesuaikan dengan API key yang valid.

Default server URL: `http://localhost:3000`

## Quick Start

1. **Import Postman Collection**: Import `postman_collection.json` ke Postman
2. **Run Manual Tests**: `./manual_test_guide.sh test`
3. **Run CLI Tester**: `go run jwt_tester.go`

Semua tools sudah disesuaikan dengan struktur endpoint dan flow authentication yang baru.
