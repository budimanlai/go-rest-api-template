package migration

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// Migrator wraps golang-migrate with our custom interface
type Migrator struct {
	migrate *migrate.Migrate
	db      *sqlx.DB
}

// NewMigrator creates a new migrator instance using golang-migrate
func NewMigrator(db *sqlx.DB, migrationsDir string) (*Migrator, error) {
	// Convert sqlx.DB to sql.DB for golang-migrate
	sqlDB := db.DB

	// Create MySQL driver instance
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create MySQL driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		"mysql",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return &Migrator{
		migrate: m,
		db:      db,
	}, nil
}

// Up applies all pending migrations
func (m *Migrator) Up() error {
	fmt.Println("🚀 Applying pending migrations...")

	err := m.migrate.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up failed: %w", err)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("✅ No pending migrations to apply")
	} else {
		fmt.Println("✅ All migrations applied successfully")
	}

	return nil
}

// Down rolls back the specified number of migrations
func (m *Migrator) Down(steps int) error {
	fmt.Printf("📉 Rolling back %d migration(s)...\n", steps)

	err := m.migrate.Steps(-steps)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration down failed: %w", err)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("✅ No migrations to rollback")
	} else {
		fmt.Printf("✅ Successfully rolled back %d migration(s)\n", steps)
	}

	return nil
}

// Status shows the current migration status
func (m *Migrator) Status() error {
	fmt.Println("📊 Migration Status:")
	fmt.Println("===================")

	version, dirty, err := m.migrate.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		fmt.Println("❌ No migrations have been applied yet")
		return nil
	}

	if dirty {
		fmt.Printf("⚠️  Current version: %d (DIRTY - migration failed)\n", version)
		fmt.Println("💡 Fix the migration issue and run 'force' command if needed")
	} else {
		fmt.Printf("✅ Current version: %d (clean)\n", version)
	}

	return nil
}

// Create generates a new migration file (simple implementation)
func (m *Migrator) Create(name string) error {
	// golang-migrate doesn't have built-in create command in library
	// We'll implement a simple version
	return fmt.Errorf("migration creation should be done using golang-migrate CLI tool:\n"+
		"Install: go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest\n"+
		"Create: migrate create -ext sql -dir migrations %s", name)
}

// Force sets the migration version without running the migration
func (m *Migrator) Force(version int) error {
	fmt.Printf("🔧 Forcing migration version to %d...\n", version)

	err := m.migrate.Force(version)
	if err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	fmt.Printf("✅ Migration version forced to %d\n", version)
	return nil
}

// GetDB returns the database connection
func (m *Migrator) GetDB() *sqlx.DB {
	return m.db
}

// Close closes the migration instance
func (m *Migrator) Close() error {
	sourceErr, databaseErr := m.migrate.Close()
	if sourceErr != nil {
		return sourceErr
	}
	return databaseErr
}
