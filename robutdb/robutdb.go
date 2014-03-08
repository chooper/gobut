
package robutdb

import (
    "log"
    "os"
    _ "github.com/lib/pq"
    "database/sql"
)

func SaveURL(url string) error {
    // Grab env var, return if not set
    database_url := os.Getenv("DATABASE_URL")
    if database_url == "" {
        return nil
    }

    // Connect to DB
    db, err := sql.Open("postgres", database_url)
    if err != nil {
        log.Print(err)
    }

    // Insert URL into DB
    _, err = db.Exec("INSERT INTO urls (\"when\", url) VALUES (NOW(), $1)", url)
    if err != nil {
        log.Print(err)
    }
    return nil
}
