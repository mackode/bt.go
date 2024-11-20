package main

import (
  "time"
)

func (tr *trader) strat_hold() tradeFu {
  return func(dt time.Time, quote float64) {
    if !tr.holds {
      tr.buy(dt, quote)
    }
  }
}

func (tr *trader) strat_buydrop() tradeFu {
  return func(dt time.Time, quote float64) {
    if tr.prevQ != 0 {
      if tr.holds {
        if quote > 1.1 * tr.cost || quote < 0.9 * tr.cost {
          tr.sell(dt, quote)
        }
      } else {
        if quote < 0.98 * tr.prevQ {
          tr.buy(dt, quote)
        }
      }
    }
  }
}

func (tr *trader) strat_firstweek() tradeFu {
  held := 0
  return func(dt time.Time, quote float64) {
    if tr.holds {
      held += 1
      if held > 5 {
        tr.sell(dt, quote)
        held = 0
      }
    } else {
      if dt.Day() < 7 {
        tr.buy(dt, quote)
      }
    }
  }
}

