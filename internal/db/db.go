package dbutils

import (
	"log"
	"sapopinguino/internal/config"

	"github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConfigDB() {
    var err error

    DB, err = sqlx.Connect(
        "postgres", 
        config.C.Database.DSN,
    )
    if err != nil {
        log.Fatalf(
            "Error  connecting to database: %s",
            err,
        )
    }
}

