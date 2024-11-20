package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

type DB struct {
    Client *sql.DB
}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

// NewDatabase - New db connection with a particular model
func NewDatabase(dsn string) (*DB, error) {
    dbConn := &DB{}

    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, err
    }
    db.SetMaxOpenConns(maxOpenDbConn)
    db.SetMaxIdleConns(maxIdleDbConn)
    db.SetConnMaxLifetime(maxDbLifeTime)

    err = checkDB(db)
    if err != nil {
        log.Fatal(err)
    }

    dbConn.Client = db

    return dbConn, nil
}

// checkDB - will check db connection
func checkDB(d *sql.DB) error {
    err := d.Ping()
    if err != nil {
        fmt.Println("Error", err)
        return err
    }
    fmt.Println("*** Pinged database successfully ***")
    return nil
}
