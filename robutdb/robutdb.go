
package robutdb

import (
    "errors"
    "log"
    "os"
    _ "github.com/lib/pq"
    "database/sql"
)

func RandomURL() (string, error) {
    // Grab env var, return if not set
    database_url := os.Getenv("DATABASE_URL")
    if database_url == "" {
        return "", errors.New("DATABASE_URL not set")
    }

    // Connect to DB
    db, err := sql.Open("postgres", database_url)

    if err != nil {
        log.Print(err)
        return "", err
    }

    // TODO This might become slow one day
    rows, err := db.Query(`SELECT url
        FROM urls
        ORDER BY RANDOM()
        LIMIT 1`)
    defer rows.Close()

    if err != nil {
        log.Print(err)
        return "", err
    }

    // Select first & only result
    rows.Next()
    var url string
    err = rows.Scan(&url)
    if err != nil {
        log.Print(err)
        return "", err
    }
    return url, nil
}

func TopSharers() (map[string]int, error) {
    // Grab env var, return if not set
    database_url := os.Getenv("DATABASE_URL")
    if database_url == "" {
        return nil, errors.New("DATABASE_URL not set")
    }

    // Connect to DB
    db, err := sql.Open("postgres", database_url)

    if err != nil {
        log.Print(err)
        return nil, err
    }

    // TODO Don't hardcode the limit
    rows, err := db.Query(`WITH urlstats AS (
            SELECT
                split_part(shared_by, '!', 1) as nick,
                COUNT(*) as count
            FROM urls
            WHERE "when" > NOW() - INTERVAL '1 week'
            GROUP BY split_part(shared_by, '!', 1)
        )
        SELECT *
        FROM urlstats
        WHERE nick != ''
        ORDER BY count DESC
        LIMIT 5`)
    defer rows.Close()

    if err != nil {
        log.Print(err)
        return nil, err
    }

    var result = make(map[string]int)
    for rows.Next() {
        var nick string
        var count int
        err = rows.Scan(&nick, &count)
        if err != nil {
            log.Print(err)
            return nil, err
        }

        result[nick] = count
    }

    return result, nil
}

func SaveURL(url string, title string, prefix string) error {
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
    _, err = db.Exec("INSERT INTO urls (\"when\", url, title, shared_by) VALUES (NOW(), $1, $2, $3)",
        url, title, prefix)
    if err != nil {
        log.Print(err)
    }
    return nil
}

