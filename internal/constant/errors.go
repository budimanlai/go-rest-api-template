package constant

// Custom Error Types
const (
	// Validation Errors
	ErrValidationUsername      = "validation_username_invalid"
	ErrValidationEmail         = "validation_email_invalid"
	ErrValidationPassword      = "validation_password_weak"
	ErrValidationPasswordMatch = "validation_password_mismatch"
	ErrValidationRequired      = "validation_field_required"
	ErrValidationLength        = "validation_field_length"
	ErrValidationFormat        = "validation_field_format"

	// User Errors
	ErrUserAlreadyExists    = "user_already_exists"
	ErrUserNotActive        = "user_not_active"
	ErrUserSuspended        = "user_suspended"
	ErrUserEmailNotVerified = "user_email_not_verified"
	ErrUserAccountLocked    = "user_account_locked"

	// Authentication Errors
	ErrAuthInvalidAPIKey      = "auth_invalid_api_key"
	ErrAuthMissingToken       = "auth_missing_token"
	ErrAuthInvalidTokenFormat = "auth_invalid_token_format"
	ErrAuthTokenMalformed     = "auth_token_malformed"

	// Authorization Errors
	ErrAuthzInsufficientPermission = "authz_insufficient_permission"
	ErrAuthzResourceNotOwned       = "authz_resource_not_owned"
	ErrAuthzOperationNotAllowed    = "authz_operation_not_allowed"

	// Database Errors
	ErrDBConnectionFailed    = "db_connection_failed"
	ErrDBQueryFailed         = "db_query_failed"
	ErrDBTransactionFailed   = "db_transaction_failed"
	ErrDBConstraintViolation = "db_constraint_violation"
	ErrDBRecordNotFound      = "db_record_not_found"
	ErrDBDuplicateEntry      = "db_duplicate_entry"

	// External Service Errors
	ErrEmailServiceUnavailable = "email_service_unavailable"
	ErrEmailSendFailed         = "email_send_failed"
	ErrSMSServiceUnavailable   = "sms_service_unavailable"
	ErrCacheServiceUnavailable = "cache_service_unavailable"

	// File Upload Errors
	ErrFileTooBig         = "file_too_big"
	ErrFileTypeNotAllowed = "file_type_not_allowed"
	ErrFileCorrupted      = "file_corrupted"
	ErrFileUploadFailed   = "file_upload_failed"

	// Rate Limiting Errors
	ErrRateLimitExceeded = "rate_limit_exceeded"
	ErrTooManyRequests   = "too_many_requests"
	ErrIPBlocked         = "ip_blocked"

	// Business Logic Errors
	ErrBusinessRuleViolation = "business_rule_violation"
	ErrInvalidOperation      = "invalid_operation"
	ErrOperationNotAllowed   = "operation_not_allowed"
	ErrResourceInUse         = "resource_in_use"
	ErrResourceLocked        = "resource_locked"
)

// Success Messages for specific operations
const (
	// User Success Messages
	MsgUserRegistered     = "user_registered"
	MsgUserActivated      = "user_activated"
	MsgUserDeactivated    = "user_deactivated"
	MsgUserSuspended      = "user_suspended"
	MsgUserUnsuspended    = "user_unsuspended"
	MsgUserProfileUpdated = "user_profile_updated"

	// Authentication Success Messages
	MsgTokenValidated   = "token_validated"
	MsgTokenRevoked     = "token_revoked"
	MsgSessionCreated   = "session_created"
	MsgSessionDestroyed = "session_destroyed"

	// Email Success Messages
	MsgEmailSent     = "email_sent"
	MsgEmailVerified = "email_verified"
	MsgEmailChanged  = "email_changed"

	// General Success Messages
	MsgOperationCompleted = "operation_completed"
	MsgResourceCreated    = "resource_created"
	MsgResourceUpdated    = "resource_updated"
	MsgResourceDeleted    = "resource_deleted"
)

// Log Event Types
const (
	LogEventUserLogin          = "user_login"
	LogEventUserLogout         = "user_logout"
	LogEventUserRegistration   = "user_registration"
	LogEventPasswordChange     = "password_change"
	LogEventPasswordReset      = "password_reset"
	LogEventUserUpdate         = "user_update"
	LogEventUserDelete         = "user_delete"
	LogEventAPIKeyUsed         = "api_key_used"
	LogEventTokenGenerated     = "token_generated"
	LogEventTokenRefreshed     = "token_refreshed"
	LogEventUnauthorizedAccess = "unauthorized_access"
	LogEventSuspiciousActivity = "suspicious_activity"
)

// Audit Fields
const (
	AuditFieldAction       = "action"
	AuditFieldResource     = "resource"
	AuditFieldResourceID   = "resource_id"
	AuditFieldUserID       = "user_id"
	AuditFieldIPAddress    = "ip_address"
	AuditFieldUserAgent    = "user_agent"
	AuditFieldTimestamp    = "timestamp"
	AuditFieldSuccess      = "success"
	AuditFieldErrorCode    = "error_code"
	AuditFieldErrorMessage = "error_message"
	AuditFieldChanges      = "changes"
)

// Cache Key Patterns
const (
	CachePatternUser          = "user:%s"
	CachePatternUserSession   = "session:user:%d"
	CachePatternAPIKey        = "apikey:%s"
	CachePatternRateLimit     = "ratelimit:%s:%s"
	CachePatternEmailVerify   = "email_verify:%s"
	CachePatternPasswordReset = "password_reset:%s"
)

// Regex Patterns
const (
	RegexUsername = `^[a-zA-Z0-9_]{3,50}$`
	RegexEmail    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	RegexPassword = `^.{8,100}$`
	RegexAPIKey   = `^[a-zA-Z0-9]{32}$`
	RegexToken    = `^[a-zA-Z0-9._-]+$`
	RegexUUID     = `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
)

// Time Durations (in seconds)
const (
	DurationTokenExpiry         = 86400  // 24 hours
	DurationRefreshTokenExpiry  = 604800 // 7 days
	DurationPasswordResetExpiry = 3600   // 1 hour
	DurationEmailVerifyExpiry   = 86400  // 24 hours
	DurationSessionTimeout      = 1800   // 30 minutes
	DurationCacheShort          = 300    // 5 minutes
	DurationCacheMedium         = 3600   // 1 hour
	DurationCacheLong           = 86400  // 24 hours
)

// Feature Flags
const (
	FeatureEmailVerification = "email_verification"
	FeatureTwoFactorAuth     = "two_factor_auth"
	FeaturePasswordHistory   = "password_history"
	FeatureAccountLocking    = "account_locking"
	FeatureAuditLogging      = "audit_logging"
	FeatureRateLimiting      = "rate_limiting"
	FeatureAPIVersioning     = "api_versioning"
)

// Environment-specific Constants
const (
	ConfigKeyDBHost            = "DB_HOST"
	ConfigKeyDBPort            = "DB_PORT"
	ConfigKeyDBName            = "DB_NAME"
	ConfigKeyDBUser            = "DB_USER"
	ConfigKeyDBPassword        = "DB_PASSWORD"
	ConfigKeyJWTSecret         = "JWT_SECRET"
	ConfigKeyAPISecret         = "API_SECRET"
	ConfigKeyEmailSMTPHost     = "EMAIL_SMTP_HOST"
	ConfigKeyEmailSMTPPort     = "EMAIL_SMTP_PORT"
	ConfigKeyEmailSMTPUser     = "EMAIL_SMTP_USER"
	ConfigKeyEmailSMTPPassword = "EMAIL_SMTP_PASSWORD"
	ConfigKeyRedisHost         = "REDIS_HOST"
	ConfigKeyRedisPort         = "REDIS_PORT"
	ConfigKeyRedisPassword     = "REDIS_PASSWORD"
)

// Notification Types
const (
	NotificationTypeEmail = "email"
	NotificationTypeSMS   = "sms"
	NotificationTypePush  = "push"
	NotificationTypeInApp = "in_app"
)

// Priority Levels
const (
	PriorityLow      = "low"
	PriorityMedium   = "medium"
	PriorityHigh     = "high"
	PriorityCritical = "critical"
)
