package repositories

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/rs/zerolog"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed migrations/*.sql
var migrations embed.FS

// TODO("Выпилить и запускать миграцию через докер")

func ApplyMigration(logger zerolog.Logger, repo Repository) error {
	logger.Debug().Msg("repository: running DB migrations")
	sub, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return err
	}

	migrationSource := &migrate.HttpFileSystemMigrationSource{FileSystem: http.FS(sub)}
	n, err := migrate.Exec(repo.SqlxDB().DB, "postgres", migrationSource, migrate.Up)
	if err != nil {
		return err
	}
	logger.Debug().Msgf("repository: applied %v migrations", n)

	return nil
}
