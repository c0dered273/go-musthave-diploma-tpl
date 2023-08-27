package repositories

import (
	"database/sql"
	"embed"
	"io/fs"
	"net/http"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed migrations/*.sql
var migrations embed.FS

// TODO("Запускать миграции через докер")

func ApplyMigration(logger zerolog.Logger, cfg *configs.ServerConfig) error {
	logger.Debug().Msg("repository: running DB migrations")
	sub, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return err
	}

	db, err := sql.Open("pgx", cfg.DatabaseURI)
	if err != nil {
		return err
	}

	migrationSource := &migrate.HttpFileSystemMigrationSource{FileSystem: http.FS(sub)}
	n, err := migrate.Exec(db, "postgres", migrationSource, migrate.Up)
	if err != nil {
		return err
	}
	logger.Debug().Msgf("repository: applied %v migrations", n)

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}
