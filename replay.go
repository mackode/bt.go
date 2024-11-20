package main

import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "time"
)

func replay(db *sql.DB, cb func(time.Time, float64)) error {
  query := `SELECT date, quote FROM quotes ORDER BY date ASC`
  rows, err := db.Query(query)
  if err != nil {
    return err
  }
  defer rows.Close()

  for rows.Next() {
    var date string
    var quote float64
    err := rows.Scan(&date, &quote)
    if err != nil {
      return err
    }
    dt, err := time.Parse("2006-01-02T15:00:05Z07:00", date)
    if err != nil {
      return err
    }
    cb(dt, quote)
  }

  if err = rows.Err(); err != nil {
    return err
  }

  return nil
}

