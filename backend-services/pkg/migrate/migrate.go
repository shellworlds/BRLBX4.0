package migrate

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Up applies migrations from a local directory to the given DSN (postgres://...).
func Up(dsn, migrationsDir string) error {
	if dsn == "" || migrationsDir == "" {
		return fmt.Errorf("migrate: missing dsn or dir")
	}
	m, err := migrate.New("file://"+migrationsDir, dsn)
	if err != nil {
		return err
	}
	defer m.Close()
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}
