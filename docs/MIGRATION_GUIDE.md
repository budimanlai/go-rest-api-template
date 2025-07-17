# 🗃️ Database Migration Guide

## 📖 Overview

The migration system uses **golang-migrate** - a production-ready database migration tool that provides versioned database schema management with rollback capabilities.

## 🏗️ Migration Structure

```
migrations/
├── 001_create_users_table.up.sql
├── 001_create_users_table.down.sql
├── 002_create_api_keys_table.up.sql
├── 002_create_api_keys_table.down.sql
└── ... (future migrations)

pkg/migration/
└── migrator.go              # golang-migrate wrapper
```

## 🛠️ **Why golang-migrate?**

✅ **Production-tested** by thousands of companies  
✅ **Battle-tested** with extensive community support  
✅ **Low-maintenance** - no custom migration code to maintain  
✅ **Multi-database** support (MySQL, PostgreSQL, SQLite, etc.)  
✅ **Atomic migrations** with transaction safety  
✅ **CLI tool** for advanced operations

## 🚀 Available Commands

### Apply Migrations
```bash
# Apply all pending migrations
./rest-api migrate-up

# Apply migrations in production
./rest-api migrate-up --env=production
```

### Rollback Migrations
```bash
# Rollback 1 migration (default)
./rest-api migrate-down

# Rollback 3 migrations
./rest-api migrate-down --steps=3
```

### Check Migration Status
```bash
# Show which migrations are applied/pending
./rest-api migrate-status
```

### Create New Migration
```bash
# Method 1: Using golang-migrate CLI (Recommended)
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate create -ext sql -dir migrations create_products_table

# Method 2: Using template command (limited)
./rest-api migrate-create --name=create_products_table
```

### Force Migration (Fix Dirty State)
```bash
# If migration fails and database is in dirty state
./rest-api migrate-force --version=3
```

### Reset All Migrations
```bash
# Rollback all migrations then apply all again
./rest-api migrate-reset
```

### Fresh Migration (⚠️ DESTRUCTIVE)
```bash
# Drop ALL tables and run migrations from scratch
./rest-api migrate-fresh

# WARNING: This will delete ALL data!
```

## 📝 Migration File Format

golang-migrate uses **separate files** for up and down migrations:

### Up Migration (001_create_products_table.up.sql)
```sql
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    category_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_category_id (category_id),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
);

-- Insert default data
INSERT INTO products (name, price) VALUES 
('Sample Product 1', 99.99),
('Sample Product 2', 149.99);
```

### Down Migration (001_create_products_table.down.sql)
```sql
DROP TABLE IF EXISTS products;
```

## 🎯 Best Practices

### 1. **Naming Convention**
```
Pattern: {version}_{description}.{direction}.sql

Examples:
- 001_create_users_table.up.sql / 001_create_users_table.down.sql
- 002_add_email_index_to_users.up.sql / 002_add_email_index_to_users.down.sql
- 003_create_products_table.up.sql / 003_create_products_table.down.sql
```

### 2. **Migration Content Guidelines**
```sql
-- ✅ GOOD: Separate up and down files
-- File: 001_create_categories.up.sql
CREATE TABLE categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

-- File: 001_create_categories.down.sql
DROP TABLE IF EXISTS categories;

-- ✅ GOOD: Add indexes for performance
-- File: 002_add_user_indexes.up.sql
CREATE INDEX idx_user_email ON users(email);
CREATE INDEX idx_user_status ON users(status);

-- File: 002_add_user_indexes.down.sql
DROP INDEX idx_user_email ON users;
DROP INDEX idx_user_status ON users;
```

### 3. **Production Safety**
```bash
# Always check status before applying
./rest-api migrate-status

# Apply migrations one by one in production
./rest-api migrate-down --steps=1  # If rollback needed
./rest-api migrate-up               # Apply pending

# Never use migrate-fresh in production!
```

## 🔧 Advanced Usage

### Custom Migration Directory
```go
// In your code, you can specify custom migration directory
migrator := migration.NewMigrator(db, "./custom-migrations")
```

### Migration with Transactions
The migrator automatically wraps each migration in a transaction, so if any SQL statement fails, the entire migration is rolled back.

### Check Applied Migrations
```go
applied, err := migrator.GetAppliedMigrations()
for _, migration := range applied {
    fmt.Printf("Applied: %s at %s\n", migration.Name, migration.AppliedAt)
}
```

## 📊 Migration Status Output

```bash
$ ./rest-api migrate-status

📊 Migration Status:
===================
✅ Current version: 3 (clean)
```

**Status Meanings:**
- **Clean**: All migrations applied successfully
- **Dirty**: Last migration failed, needs manual intervention
- **No migrations**: Database is empty, ready for first migration

## 🛠️ Database Schema Tracking

golang-migrate creates a `schema_migrations` table to track applied migrations:

```sql
CREATE TABLE schema_migrations (
    version bigint NOT NULL PRIMARY KEY,
    dirty boolean NOT NULL
);
```

**Simple & Efficient**: Only tracks version and dirty state.

## 🚨 Common Pitfalls & Solutions

### 1. **Failed Migration Recovery**
```bash
# If a migration fails, check the status
./rest-api migrate-status

# If database is in dirty state, force to last good version
./rest-api migrate-force --version=2

# Then fix the migration file and try again
./rest-api migrate-up
```

### 2. **Rollback Limitations**
```sql
-- ✅ GOOD: Reversible operations
-- Up: 001_add_column.up.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- Down: 001_add_column.down.sql  
ALTER TABLE users DROP COLUMN phone;

-- ❌ BAD: Data loss operations
-- Up: 002_delete_old_data.up.sql
DELETE FROM users WHERE created_at < '2023-01-01';

-- Down: 002_delete_old_data.down.sql
-- Cannot restore deleted data!
```

### 3. **Index Creation on Large Tables**
```sql
-- ✅ GOOD: Add indexes concurrently for large tables
-- Note: MySQL doesn't support CONCURRENTLY like PostgreSQL
-- For large tables, consider maintenance windows

CREATE INDEX idx_users_email ON users(email);
```

## 🔐 Security Considerations

1. **Database Permissions**: Ensure migration user has appropriate DDL permissions
2. **Backup Before Migration**: Always backup production database before major schema changes
3. **Review Process**: All migration files should be code-reviewed
4. **Testing**: Test migrations on staging environment first

## 📈 Monitoring & Logging

The golang-migrate wrapper provides comprehensive logging:
- Each migration execution with status feedback
- Clear error messages for failed migrations  
- Simple status reporting (clean/dirty state)
- **Low-maintenance**: Leverages battle-tested golang-migrate

## 🎉 Production Deployment Workflow

```bash
# 1. Create migration in development  
migrate create -ext sql -dir migrations add_new_feature

# 2. Test locally
./rest-api migrate-up
./rest-api migrate-down --steps=1  # Test rollback
./rest-api migrate-up               # Re-apply

# 3. Deploy to staging
git push staging
./rest-api migrate-up

# 4. Deploy to production (with backup)
mysqldump database_name > backup_$(date +%Y%m%d_%H%M%S).sql
./rest-api migrate-up

# 5. Verify
./rest-api migrate-status
```

---

## 🎯 **Summary: Why golang-migrate?**

✅ **Production-Ready**: Used by thousands of companies  
✅ **Low-Maintenance**: No custom migration code to maintain  
✅ **Battle-Tested**: Extensive community support & testing  
✅ **Simple API**: Clean, minimal interface  
✅ **Reliable**: Atomic transactions and dirty state handling  

This migration system ensures **safe, versioned, and low-maintenance** database schema management! 🚀
