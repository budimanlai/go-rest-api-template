#!/bin/bash

# Update UserHandler to use i18n responses
# This script replaces old response calls with i18n versions

cd /Users/budiman/Documents/development/my_github/go-rest_api

echo "🔄 Updating UserHandler to use i18n responses..."

# Make backup first
cp internal/handler/user_handler.go internal/handler/user_handler.go.backup

echo "📁 Backup created: user_handler.go.backup"

# Replace all old response calls with i18n versions
sed -i '' \
    -e 's/return response\.InternalServerError(c, "Failed to retrieve users", err\.Error())/return h.responseHelper.ErrorWithI18n(c, 500, "failed_retrieve_users", nil)/g' \
    -e 's/return response\.InternalServerError(c, "Failed to get user count", err\.Error())/return h.responseHelper.ErrorWithI18n(c, 500, "failed_get_user_count", nil)/g' \
    -e 's/return response\.Success(c, "Users retrieved successfully", paginationData)/return h.responseHelper.SuccessWithI18n(c, "users_retrieved", paginationData, nil)/g' \
    -e 's/return response\.BadRequest(c, "Invalid user ID", err\.Error())/return h.responseHelper.ErrorWithI18n(c, 400, "invalid_user_id", nil)/g' \
    -e 's/return response\.BadRequest(c, "Invalid request body", err\.Error())/return h.responseHelper.ErrorWithI18n(c, 400, "invalid_request_body", nil)/g' \
    -e 's/return response\.BadRequest(c, "Validation failed", err\.Error())/return h.responseHelper.ErrorWithI18n(c, 400, "validation_failed", nil)/g' \
    -e 's/return response\.BadRequest(c, "Password hashing failed", err\.Error())/return h.responseHelper.ErrorWithI18n(c, 400, "password_hashing_failed", nil)/g' \
    -e 's/return response\.InternalServerError(c, "Failed to update user", err\.Error())/return h.responseHelper.ErrorWithI18n(c, 500, "failed_update_user", nil)/g' \
    -e 's/return response\.InternalServerError(c, "Failed to retrieve updated user", err\.Error())/return h.responseHelper.ErrorWithI18n(c, 500, "failed_retrieve_updated_user", nil)/g' \
    -e 's/return response\.Success(c, "User updated successfully", userResponse)/return h.responseHelper.SuccessWithI18n(c, "user_updated", userResponse, nil)/g' \
    -e 's/return response\.InternalServerError(c, "Failed to delete user", err\.Error())/return h.responseHelper.ErrorWithI18n(c, 500, "failed_delete_user", nil)/g' \
    -e 's/return response\.Success(c, "User deleted successfully", nil)/return h.responseHelper.SuccessWithI18n(c, "user_deleted", nil, nil)/g' \
    -e 's/return response\.InternalServerError(c, "Failed to process forgot password", err\.Error())/return h.responseHelper.ErrorWithI18n(c, 500, "failed_forgot_password", nil)/g' \
    -e 's/return response\.Success(c, "Reset password instructions sent to your email", nil)/return h.responseHelper.SuccessWithI18n(c, "reset_password_sent", nil, nil)/g' \
    -e 's/return response\.BadRequest(c, "Reset password failed", err\.Error())/return h.responseHelper.ErrorWithI18n(c, 400, "reset_password_failed", nil)/g' \
    internal/handler/user_handler.go

echo "✅ UserHandler updated successfully!"

# Test compilation
echo "🧪 Testing compilation..."
if go build ./cmd/api/; then
    echo "✅ Compilation successful!"
    rm -f api  # Remove generated binary
else
    echo "❌ Compilation failed! Restoring backup..."
    mv internal/handler/user_handler.go.backup internal/handler/user_handler.go
    echo "🔄 Backup restored."
    exit 1
fi

echo "🎉 Update completed successfully!"
echo "📝 All UserHandler methods now use i18n responses."
