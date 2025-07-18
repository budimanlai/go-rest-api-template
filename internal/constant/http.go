package constant

import "net/http"

// HTTP Status Codes (extending standard library)
const (
	StatusOK                  = http.StatusOK                  // 200
	StatusCreated             = http.StatusCreated             // 201
	StatusAccepted            = http.StatusAccepted            // 202
	StatusNoContent           = http.StatusNoContent           // 204
	StatusBadRequest          = http.StatusBadRequest          // 400
	StatusUnauthorized        = http.StatusUnauthorized        // 401
	StatusForbidden           = http.StatusForbidden           // 403
	StatusNotFound            = http.StatusNotFound            // 404
	StatusMethodNotAllowed    = http.StatusMethodNotAllowed    // 405
	StatusConflict            = http.StatusConflict            // 409
	StatusUnprocessableEntity = http.StatusUnprocessableEntity // 422
	StatusTooManyRequests     = http.StatusTooManyRequests     // 429
	StatusInternalServerError = http.StatusInternalServerError // 500
	StatusBadGateway          = http.StatusBadGateway          // 502
	StatusServiceUnavailable  = http.StatusServiceUnavailable  // 503
)

// HTTP Methods
const (
	MethodGET     = http.MethodGet
	MethodPOST    = http.MethodPost
	MethodPUT     = http.MethodPut
	MethodPATCH   = http.MethodPatch
	MethodDELETE  = http.MethodDelete
	MethodOPTIONS = http.MethodOptions
	MethodHEAD    = http.MethodHead
)

// API Route Patterns
const (
	// Auth Routes
	RouteAuthRegister       = "/auth/register"
	RouteAuthLogin          = "/auth/login"
	RouteAuthLogout         = "/auth/logout"
	RouteAuthRefresh        = "/auth/refresh"
	RouteAuthVerify         = "/auth/verify"
	RouteAuthForgot         = "/auth/forgot-password"
	RouteAuthReset          = "/auth/reset-password"
	RouteAuthChangePassword = "/auth/change-password"

	// User Routes
	RouteUsers           = "/users"
	RouteUserByID        = "/users/:id"
	RouteUserProfile     = "/users/profile"
	RouteUserDynamic     = "/users/dynamic"
	RouteUserDynamicByID = "/users/dynamic/:id"
	RouteUserSafe        = "/users/safe"

	// Health Check Routes
	RouteHealth          = "/health"
	RouteHealthLiveness  = "/health/liveness"
	RouteHealthReadiness = "/health/readiness"

	// API Documentation Routes
	RouteSwagger    = "/swagger/*"
	RouteAPIVersion = "/version"
)

// CORS Constants
const (
	CORSAllowOrigins     = "*"
	CORSAllowMethods     = "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS"
	CORSAllowHeaders     = "Origin,Content-Length,Content-Type,Authorization,X-API-Key"
	CORSExposeHeaders    = "Content-Length,Content-Range"
	CORSMaxAge           = 86400 // 24 hours
	CORSAllowCredentials = true
)

// Request/Response Constants
const (
	// Query Parameters
	QueryParamPage   = "page"
	QueryParamLimit  = "limit"
	QueryParamOffset = "offset"
	QueryParamSort   = "sort"
	QueryParamOrder  = "order"
	QueryParamSearch = "search"
	QueryParamFilter = "filter"
	QueryParamFields = "fields"

	// Sort Orders
	SortOrderASC  = "asc"
	SortOrderDESC = "desc"

	// Default Sort Field
	DefaultSortField = "created_at"
)

// Middleware Constants
const (
	MiddlewareTimeout     = 30       // seconds
	MiddlewareMaxBodySize = 10 << 20 // 10MB
)

// Response Format Constants
const (
	ResponseFieldSuccess    = "success"
	ResponseFieldMessage    = "message"
	ResponseFieldData       = "data"
	ResponseFieldError      = "error"
	ResponseFieldPagination = "pagination"
	ResponseFieldMeta       = "meta"
)

// File Upload Response Constants
const (
	UploadFieldFile    = "file"
	UploadFieldFiles   = "files"
	UploadMaxMemory    = 10 << 20 // 10MB
	UploadTempDir      = "/tmp"
	UploadAllowedTypes = "image/jpeg,image/png,image/gif,application/pdf"
)

// Security Headers
const (
	HeaderXFrameOptions         = "X-Frame-Options"
	HeaderXContentTypeOptions   = "X-Content-Type-Options"
	HeaderXXSSProtection        = "X-XSS-Protection"
	HeaderStrictTransportSec    = "Strict-Transport-Security"
	HeaderContentSecurityPolicy = "Content-Security-Policy"
	HeaderReferrerPolicy        = "Referrer-Policy"

	// Security Header Values
	ValueXFrameOptionsDeny    = "DENY"
	ValueXContentTypeNosniff  = "nosniff"
	ValueXXSSProtectionBlock  = "1; mode=block"
	ValueStrictTransportSec   = "max-age=31536000; includeSubDomains"
	ValueReferrerPolicyStrict = "strict-origin-when-cross-origin"
)

// Request Context Keys
const (
	ContextKeyRequestID = "request_id"
	ContextKeyStartTime = "start_time"
	ContextKeyUserAgent = "user_agent"
	ContextKeyIPAddress = "ip_address"
	ContextKeyRealIP    = "real_ip"
)

// API Versioning
const (
	HeaderAPIVersion    = "X-API-Version"
	HeaderAcceptVersion = "Accept-Version"
	CurrentAPIVersion   = "1.0.0"
)

// Rate Limiting Headers
const (
	HeaderRateLimitLimit     = "X-RateLimit-Limit"
	HeaderRateLimitRemaining = "X-RateLimit-Remaining"
	HeaderRateLimitReset     = "X-RateLimit-Reset"
	HeaderRetryAfter         = "Retry-After"
)
