// TODO: implement postgres storage provider
package postgresdb

import (
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseApp struct {
	log *slog.Logger
	dsn string
	db  *gorm.DB
}

// New create DatabaseApp
func New(log *slog.Logger, dsn string) *DatabaseApp {
	return &DatabaseApp{
		log: log,
		dsn: dsn,
	}
}

// MustRun create postgres database connection. Panic if failed
func (d *DatabaseApp) MustRun() {
	err := d.Run()
	if err != nil {
		panic(err)
	}
}

// Run create postgres database connection
func (d *DatabaseApp) Run() error {
	db, err := gorm.Open(postgres.Open(d.dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Database connect error")
	}
	d.db = db
	return nil
}

// Stop close postgres database connection
func (d *DatabaseApp) Stop() error {
	if d.db == nil {
		return fmt.Errorf("Database is not connect")
	}

	conn, err := d.db.DB()

	if err != nil {
		return fmt.Errorf("Failed to get conection")
	}

	err = conn.Close()
	if err != nil {
		return fmt.Errorf("Failed to close")
	}

	return nil
}
