package constant

// Database Table Names
const (
	TableUser = "user"
)

// User Status Constants
const (
	UserStatusActive    = "active"
	UserStatusInactive  = "inactive"
	UserStatusSuspended = "suspended"
)

// Default Values
const (
	DefaultCreatedBy = 0 // System user
	DefaultUpdatedBy = 0 // System user
	AuthKeyLength    = 32
)

// Pagination Constants
const (
	DefaultPageLimit  = 10
	MaxPageLimit      = 100
	DefaultPageOffset = 0
)

// Validation Constants
const (
	MinUsernameLength = 3
	MaxUsernameLength = 50
	MinPasswordLength = 8
	MaxPasswordLength = 100
	MaxEmailLength    = 100
)

// HTTP Status Messages
const (
	MsgUserCreated       = "user_created"
	MsgUserRetrieved     = "user_retrieved"
	MsgUsersRetrieved    = "users_retrieved"
	MsgUserUpdated       = "user_updated"
	MsgUserDeleted       = "user_deleted"
	MsgUserVerified      = "user_verified"
	MsgPasswordChanged   = "password_changed"
	MsgPasswordResetSent = "password_reset_sent"
	MsgPasswordReset     = "password_reset"
	MsgLoginSuccessful   = "login_successful"
	MsgLogoutSuccessful  = "logout_successful"
	MsgTokenGenerated    = "token_generated"
	MsgTokenRefreshed    = "token_refreshed"
)

// Error Messages
const (
	ErrInvalidRequest        = "invalid_request"
	ErrValidationFailed      = "validation_failed"
	ErrUserNotFound          = "user_not_found"
	ErrUsernameExists        = "username_exists"
	ErrEmailExists           = "email_exists"
	ErrInvalidCredentials    = "invalid_credentials"
	ErrInvalidUserID         = "invalid_user_id"
	ErrInternalServer        = "internal_server"
	ErrPasswordHashingFailed = "password_hashing_failed"
	ErrFailedUpdateUser      = "failed_update_user"
	ErrFailedRetrieveUser    = "failed_retrieve_updated_user"
	ErrInvalidToken          = "invalid_token"
	ErrTokenExpired          = "token_expired"
	ErrUnauthorized          = "unauthorized"
	ErrForbidden             = "forbidden"
)

// JWT Constants
const (
	JWTIssuer          = "go-rest-api-template"
	JWTAudience        = "api-users"
	DefaultJWTExpiry   = 24 * 60 * 60     // 24 hours in seconds
	RefreshTokenExpiry = 7 * 24 * 60 * 60 // 7 days in seconds
)

// API Constants
const (
	APIVersion = "v1"
	APIPrefix  = "/api/" + APIVersion
)

// Context Keys
const (
	ContextKeyUserID   = "user_id"
	ContextKeyUsername = "username"
	ContextKeyAPIKey   = "api_key"
)

// Header Constants
const (
	HeaderAPIKey        = "X-API-Key"
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"
	HeaderAccept        = "Accept"
)

// Content Types
const (
	ContentTypeJSON = "application/json"
	ContentTypeXML  = "application/xml"
	ContentTypeForm = "application/x-www-form-urlencoded"
)

// Environment Constants
const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
	EnvTesting     = "testing"
)

// Logging Constants
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelFatal = "fatal"
)

// Cache Constants
const (
	CacheKeyUserByID       = "user:id:%d"
	CacheKeyUserByEmail    = "user:email:%s"
	CacheKeyUserByUsername = "user:username:%s"
	DefaultCacheTTL        = 300 // 5 minutes in seconds
)

// File Upload Constants
const (
	MaxFileSize     = 10 << 20 // 10MB
	AllowedImageExt = ".jpg,.jpeg,.png,.gif"
	AllowedDocExt   = ".pdf,.doc,.docx,.txt"
	UploadPath      = "./uploads"
)

// Rate Limiting Constants
const (
	DefaultRateLimit       = 100 // requests per minute
	AuthRateLimit          = 10  // auth requests per minute
	PasswordResetRateLimit = 5   // password reset per hour
)

// Email Templates
const (
	EmailTemplateVerification  = "verification"
	EmailTemplatePasswordReset = "password_reset"
	EmailTemplateWelcome       = "welcome"
)
