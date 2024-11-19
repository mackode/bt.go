package main

import (
  "flag"
  "fmt"
  "log"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func main() {
  update := flag.Bool("update", false, "update quotes in db")
  strategy := flag.String("strategy", "hold", "trader strategy")
  flag.Parse()

  db, err := sql.Open("sqlite3", "./quotes.db")
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  if *update {
    err := updater(db, "nflx")
    if err != nil {
      log.Fatal(err)
    }
    return
  }

  tr := newTrader(*strategy)
  err = replay(db, tr.trade)
  if err != nil {
    log.Fatal(err)
  }

  if tr.holds {
    tr.sell(tr.prevDt, tr.prevQ)
  }
  fmt.Printf("Suma: %+.2f\n", tr.ledger)
}
