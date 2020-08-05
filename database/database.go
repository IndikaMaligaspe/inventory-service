package database

import (
	"database/sql"
	"log"
	"time"
)

// DbConn : Connection object
var DbConn *sql.DB

// SetupDatabase : Connection to database objects
func SetupDatabase() {
	var err error
	DbConn, err = sql.Open("mysql", "root:kissme@tcp(127.0.0.1:3306)/inventorydb")
	DbConn.SetConnMaxLifetime(60 * time.Second)
	DbConn.SetMaxIdleConns(4)
	DbConn.SetMaxOpenConns(4)
	if err != nil {
		log.Fatal(err)
	}
}
