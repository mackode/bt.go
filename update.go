package main

import (
  "database/sql"
  "io/ioutil"
  "net/http"
  "net/url"
  _ "github.com/mattn/go-sqlite3"
  "github.com/tidwall/gjson"
)

func updater(db *sql.DB, ticker string) error {
  createTableSQL := `CREATE TABLE IF NOT EXISTS quotes (
    "date" DATE NOT NULL,
    "quote" REAL NOT NULL,
    UNIQUE(date, quote)
  );`
  _, err := db.Exec(createTableSQL)
  if err != nil {
    return err
  }
  
  dates, quotes, err := fetchQ(ticker)
  if err != nil {
    return err
  }

  insertSQL := `INSERT OR REPLACE INTO quotes(date, quote) VALUES (?, ?)`
  statement, err := db.Prepare(insertSQL)
  if err != nil {
    return err
  }
  defer statement.Close()

  for i, date := range dates {
    quote := quotes[i]
    _, err := statement.Exec(date.String(), quote.String())
    if err != nil {
      return err
    }
  }

  return nil
}

func fetchQ(symbols string) ([]gjson.Result, []gjson.Result, error) {
  dates := []gjson.Result{}
  quotes := []gjson.Result{}
  u := url.URL{
    Scheme: "https",
    Host: "api.twelvedata.com",
    Path: "time_series",
  }

  q := u.Query()
  q.Set("symbol", symbols)
  q.Set("interval", "1day")
  q.Set("start_date", "2024-01-01")
  q.Set("apikey", "XYZ")
  u.RawQuery = q.Encode()
  resp, err := http.Get(u.String())
  if err != nil {
    return dates, quotes, err
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return dates, quotes, err
  }
  dates = gjson.Get(string(body), "values.#.datetime").Array()
  quotes = gjson.Get(string(body), "values.#.close").Array()

  return dates, quotes, nil
}

