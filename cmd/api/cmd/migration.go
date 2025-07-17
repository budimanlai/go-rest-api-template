package cmd

import (
	"go-rest-api-template/internal/application"

	gocli "github.com/budimanlai/go-cli"
)

// RegisterMigrationCommands registers all migration-related commands
func RegisterMigrationCommands(cli *gocli.Cli) {
	// Register migration up command
	cli.AddCommand("migrate-up", application.MigrateUpService)

	// Register migration down command
	cli.AddCommand("migrate-down", application.MigrateDownService)

	// Register migration status command
	cli.AddCommand("migrate-status", application.MigrateStatusService)

	// Register migration create command
	cli.AddCommand("migrate-create", application.MigrateCreateService)

	// Register migration reset command
	cli.AddCommand("migrate-reset", application.MigrateResetService)

	// Register migration fresh command
	cli.AddCommand("migrate-fresh", application.MigrateFreshService)

	// Register migration force command (for fixing dirty state)
	cli.AddCommand("migrate-force", application.MigrateForceService)
}
