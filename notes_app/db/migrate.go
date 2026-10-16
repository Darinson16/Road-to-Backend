package db

import (
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
)

func DSN() string {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = os.Getenv("DBNET")
	cfg.Addr = os.Getenv("DBADDRESS")
	cfg.DBName = os.Getenv("DBNAME")
	cfg.ParseTime = true
	return cfg.FormatDSN()
}

func MigrationDir() string {
	if d := os.Getenv("MIGRATION_DIR"); d != "" {
		return d
	}
	return "db/migrations"
}

// MigrateUp applies pending migrations. No-op when RUN_MIGRATIONS=false.
func MigrateUp(dir string) error {
	if os.Getenv("RUN_MIGRATIONS") == "false" {
		return nil
	}
	sqlDB, err := sql.Open("mysql", DSN())
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}
	return goose.Up(sqlDB, dir)
}
