package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func Connection() *sql.DB {
	var db *sql.DB
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = os.Getenv("DBNET")
	cfg.Addr = os.Getenv("DBADDRESS")
	cfg.DBName = os.Getenv("DBNAME")
	cfg.ParseTime = true

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Printf("Connected to the database %v. On port: %v \n", cfg.DBName, cfg.Addr)

	return db
}
