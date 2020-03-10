package database

import (
	nativeerrors "errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/kaspanet/faucet/config"
	"github.com/pkg/errors"
	"os"
)

// db is the faucet database.
var db *pg.DB

// DB returns a reference to the database connection
func DB() (*pg.DB, error) {
	if db == nil {
		return nil, errors.New("Database is not connected")
	}
	return db, nil
}

// Connect connects to the database mentioned in the config variable.
func Connect(cfg *config.Config) error {
	migrator, driver, err := openMigrator(cfg)
	if err != nil {
		return err
	}
	isCurrent, version, err := isCurrent(migrator, driver)
	if err != nil {
		return errors.Errorf("Error checking whether the database is current: %s", err)
	}
	if !isCurrent {
		return errors.Errorf("Database is not current (version %d). Please migrate"+
			" the database by running the faucet with --migrate flag and then run it again.", version)
	}
	connectionOptions, err := pg.ParseURL(buildConnectionString(cfg))
	if err != nil {
		return err
	}

	db = pg.Connect(connectionOptions)

	return nil
}

// Close closes the connection to the database
func Close() error {
	if db == nil {
		return nil
	}
	err := db.Close()
	db = nil
	return err
}

func buildConnectionString(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)
}

// isCurrent resolves whether the database is on the latest
// version of the schema.
func isCurrent(migrator *migrate.Migrate, driver source.Driver) (bool, uint, error) {
	// Get the current version
	version, isDirty, err := migrator.Version()
	if nativeerrors.Is(err, migrate.ErrNilVersion) {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	if isDirty {
		return false, 0, errors.Errorf("Database is dirty")
	}

	// The database is current if Next returns ErrNotExist
	_, err = driver.Next(version)
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		if pathErr.Err == os.ErrNotExist {
			return true, version, nil
		}
	}
	return false, version, err
}

func openMigrator(cfg *config.Config) (*migrate.Migrate, source.Driver, error) {
	driver, err := source.Open("file://migrations")
	if err != nil {
		return nil, nil, err
	}
	migrator, err := migrate.NewWithSourceInstance(
		"migrations", driver, buildConnectionString(cfg))
	if err != nil {
		return nil, nil, err
	}
	return migrator, driver, nil
}

// Migrate database to the latest version.
func Migrate(cfg *config.Config) error {
	migrator, driver, err := openMigrator(cfg)
	if err != nil {
		return err
	}
	isCurrent, version, err := isCurrent(migrator, driver)
	if err != nil {
		return errors.Errorf("Error checking whether the database is current: %s", err)
	}
	if isCurrent {
		log.Infof("Database is already up-to-date (version %d)", version)
		return nil
	}
	err = migrator.Up()
	if err != nil {
		return err
	}
	version, isDirty, err := migrator.Version()
	if err != nil {
		return err
	}
	if isDirty {
		return errors.Errorf("error migrating database: database is dirty")
	}
	log.Infof("Migrated database to the latest version (version %d)", version)
	return nil
}
