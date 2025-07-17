#!/bin/bash

# Update repository mappings to include DeletedAt field
cd /Users/budiman/Documents/development/my_github/go-rest_api

echo "🔄 Adding DeletedAt field to all entity mappings..."

# Make backup
cp internal/repository/user_repository_impl.go internal/repository/user_repository_impl.go.backup

# Replace all entity mappings to include DeletedAt
sed -i '' \
    -e 's/UpdatedAt:              userModel\.UpdatedAt,$/UpdatedAt:              userModel.UpdatedAt,\
		DeletedAt:              userModel.DeletedAt,/g' \
    internal/repository/user_repository_impl.go

echo "✅ All entity mappings updated!"

# Test compilation
echo "🧪 Testing compilation..."
if go build ./cmd/api/; then
    echo "✅ Compilation successful!"
    rm -f api
else
    echo "❌ Compilation failed! Restoring backup..."
    mv internal/repository/user_repository_impl.go.backup internal/repository/user_repository_impl.go
    echo "🔄 Backup restored."
    exit 1
fi

echo "🎉 Repository mapping update completed!"
