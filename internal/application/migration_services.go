package application

import (
	"fmt"
	"go-rest-api-template/pkg/database"
	"go-rest-api-template/pkg/migration"
	"strconv"

	gocli "github.com/budimanlai/go-cli"
)

// MigrateUpService applies all pending migrations
func MigrateUpService(c *gocli.Cli) {
	c.Log("Running migration up...")

	migrator, err := createMigrator(c)
	if err != nil {
		c.Log(fmt.Sprintf("Failed to create migrator: %v", err))
		return
	}

	if err := migrator.Up(); err != nil {
		c.Log(fmt.Sprintf("Migration up failed: %v", err))
		return
	}

	c.Log("Migration up completed successfully!")
}

// MigrateDownService rolls back migrations
func MigrateDownService(c *gocli.Cli) {
	c.Log("Running migration down...")

	migrator, err := createMigrator(c)
	if err != nil {
		c.Log(fmt.Sprintf("Failed to create migrator: %v", err))
		return
	}

	// Get steps parameter
	stepsStr := c.Args.GetString("steps")
	if stepsStr == "" {
		stepsStr = "1" // Default to 1 step
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		c.Log(fmt.Sprintf("Invalid steps value: %s", stepsStr))
		return
	}

	if err := migrator.Down(steps); err != nil {
		c.Log(fmt.Sprintf("Migration down failed: %v", err))
		return
	}

	c.Log("Migration down completed successfully!")
}

// MigrateStatusService shows migration status
func MigrateStatusService(c *gocli.Cli) {
	c.Log("Checking migration status...")

	migrator, err := createMigrator(c)
	if err != nil {
		c.Log(fmt.Sprintf("Failed to create migrator: %v", err))
		return
	}

	if err := migrator.Status(); err != nil {
		c.Log(fmt.Sprintf("Migration status check failed: %v", err))
		return
	}
}

// MigrateCreateService creates a new migration file
func MigrateCreateService(c *gocli.Cli) {
	name := c.Args.GetString("name")
	if name == "" {
		c.Log("Migration name is required. Example: --name=create_products_table")
		c.Log("Or use golang-migrate CLI tool:")
		c.Log("  go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest")
		c.Log("  migrate create -ext sql -dir migrations create_products_table")
		return
	}

	c.Log(fmt.Sprintf("Creating migration: %s", name))

	migrator, err := createMigrator(c)
	if err != nil {
		c.Log(fmt.Sprintf("Failed to create migrator: %v", err))
		return
	}

	if err := migrator.Create(name); err != nil {
		c.Log(fmt.Sprintf("Migration create failed: %v", err))
		c.Log("Consider using golang-migrate CLI tool for better migration creation")
		return
	}

	c.Log("Migration file created successfully!")
}

// MigrateResetService resets all migrations (down all + up all)
func MigrateResetService(c *gocli.Cli) {
	c.Log("Resetting all migrations...")

	migrator, err := createMigrator(c)
	if err != nil {
		c.Log(fmt.Sprintf("Failed to create migrator: %v", err))
		return
	}

	fmt.Println("🔄 Resetting all migrations...")

	// For golang-migrate, we'll drop and recreate
	// Note: This is a simplified reset - in production you might want different strategy

	fmt.Println("📉 Dropping all migrations...")
	if err := migrator.Down(-1); err != nil { // -1 means all migrations
		c.Log(fmt.Sprintf("Migration rollback failed: %v", err))
		return
	}

	// Apply all
	fmt.Println("📈 Applying all migrations...")
	if err := migrator.Up(); err != nil {
		c.Log(fmt.Sprintf("Migration up failed: %v", err))
		return
	}

	c.Log("Migration reset completed successfully!")
}

// MigrateFreshService drops all tables and runs migrations
func MigrateFreshService(c *gocli.Cli) {
	c.Log("Running migration fresh...")

	fmt.Println("⚠️  WARNING: This will drop ALL tables in the database!")
	fmt.Print("Are you sure? (y/N): ")

	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		fmt.Println("❌ Migration fresh cancelled")
		return
	}

	migrator, err := createMigrator(c)
	if err != nil {
		c.Log(fmt.Sprintf("Failed to create migrator: %v", err))
		return
	}

	// Drop all tables
	fmt.Println("🗑️  Dropping all tables...")

	db := migrator.GetDB()
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		c.Log(fmt.Sprintf("Failed to disable foreign key checks: %v", err))
		return
	}

	// Get all tables
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		c.Log(fmt.Sprintf("Failed to get tables: %v", err))
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			c.Log(fmt.Sprintf("Failed to scan table name: %v", err))
			return
		}
		tables = append(tables, table)
	}

	// Drop each table
	for _, table := range tables {
		fmt.Printf("🗑️  Dropping table: %s\n", table)
		_, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			c.Log(fmt.Sprintf("Failed to drop table %s: %v", table, err))
			return
		}
	}

	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		c.Log(fmt.Sprintf("Failed to enable foreign key checks: %v", err))
		return
	}

	fmt.Println("✅ All tables dropped")

	// Run migrations
	fmt.Println("📈 Running fresh migrations...")
	if err := migrator.Up(); err != nil {
		c.Log(fmt.Sprintf("Migration up failed: %v", err))
		return
	}

	c.Log("Migration fresh completed successfully!")
}

// MigrateForceService forces migration to specific version (for fixing dirty state)
func MigrateForceService(c *gocli.Cli) {
	versionStr := c.Args.GetString("version")
	if versionStr == "" {
		c.Log("Migration version is required. Example: --version=3")
		return
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.Log(fmt.Sprintf("Invalid version value: %s", versionStr))
		return
	}

	c.Log(fmt.Sprintf("Forcing migration to version: %d", version))

	migrator, err := createMigrator(c)
	if err != nil {
		c.Log(fmt.Sprintf("Failed to create migrator: %v", err))
		return
	}

	if err := migrator.Force(version); err != nil {
		c.Log(fmt.Sprintf("Migration force failed: %v", err))
		return
	}

	c.Log("Migration force completed successfully!")
}

// createMigrator creates a new migrator instance with database connection
func createMigrator(c *gocli.Cli) (*migration.Migrator, error) {
	// Load configuration
	c.LoadConfig()

	// Connect to database
	db, err := database.NewConnection(database.Config{
		Host:     c.Config.GetString("database.hostname"),
		Port:     c.Config.GetString("database.port"),
		Username: c.Config.GetString("database.username"),
		Password: c.Config.GetString("database.password"),
		Database: c.Config.GetString("database.database"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create migrator
	migrator, err := migration.NewMigrator(db, "./migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return migrator, nil
}
