package calc

import (
	"../log"
	"../tools"
	"fmt"
	"time"
)

type Position struct {
  pre_open bool
  open bool
  waiting bool
  waiting_ct int
  buy bool
  open_time time.Time
  open_value float64
  close_time time.Time
  close_value float64
  diff_value float64
}

const waiting_min = 1

func Calc_potential(bids []tools.Bid) {

	var gain float64

	var win_pos, lost_pos int
  var win_vente_pos, lost_vente_pos int
  var win_achat_pos, lost_achat_pos int

	var tot_win, tot_lost float64

	fmt.Println("")
	fmt.Println("#################################################")
	log.Info("Calcul des bénéfices potentiels :")
	fmt.Println("")
	fmt.Println("")

  var pos Position
  var poss []Position
  var last_bid tools.Bid

	for _, bid := range bids {

    last_bid = bid

    if !pos.pre_open && !pos.open {
      if bid.Macd_signal > bid.Macd_absol_trigger_signal {
        fmt.Println("Pre open hausse ->", bid.Bid_at.Format("2006-01-02 15:04:05"), "for", bid.Last_bid)
        pos.buy = true
        pos.pre_open = true
        pos.waiting_ct++
      } else if bid.Macd_signal < -bid.Macd_absol_trigger_signal {
        fmt.Println("Pre open baisse ->", bid.Bid_at.Format("2006-01-02 15:04:05"), "for", bid.Last_bid)
        pos.buy = false
        pos.pre_open = true
        pos.waiting_ct++
      }
      continue
    }

    if pos.pre_open && !pos.open {
      if pos.buy && bid.Macd_signal > bid.Macd_absol_trigger_signal {
        pos.waiting_ct++
      } else if pos.buy && bid.Macd_signal <= bid.Macd_absol_trigger_signal {
        fmt.Println("Close pre open hausse")
        fmt.Println("#################################################")
        pos = Position{}
      }

      if !pos.buy && bid.Macd_signal < -bid.Macd_absol_trigger_signal {
        pos.waiting_ct++
      } else if !pos.buy && bid.Macd_signal >= -bid.Macd_absol_trigger_signal {
        fmt.Println("Close pre open baisse")
        fmt.Println("#################################################")
        pos = Position{}
      }
    }

    if pos.waiting_ct < waiting_min {
      continue
    }

    if pos.pre_open && !pos.open {
      if bid.Macd_signal > bid.Macd_absol_trigger_signal {
        pos.open_time = bid.Bid_at
        pos.open_value = bid.Last_bid
        pos.open = true
        fmt.Println("Open hausse ->", pos.open_time.Format("2006-01-02 15:04:05"), "on", pos.open_value)
      } else if bid.Macd_signal < -bid.Macd_absol_trigger_signal {
        pos.open_time = bid.Bid_at
        pos.open_value = bid.Last_bid
        pos.open = true
				fmt.Println("Open baisse ->", pos.open_time.Format("2006-01-02 15:04:05"), "on", pos.open_value)
      }
      continue
    }

    if pos.open && pos.buy {
      if bid.Macd_signal <= bid.Macd_absol_trigger_signal {
        pos.close_time = bid.Bid_at
        pos.close_value = bid.Last_bid
        pos.diff_value = pos.close_value - pos.open_value
        gain += pos.diff_value

        fmt.Println("Close hausse ->", pos.close_time.Format("2006-01-02 15:04:05"), "on", pos.close_value, "with open value :", pos.open_value, "at", pos.open_time.Format("2006-01-02 15:04:05"), "gain de :", pos.diff_value)
        fmt.Println("#################################################")

        if pos.diff_value > 0 {
          win_pos++
          win_achat_pos++
          tot_win += pos.diff_value
        } else {
          lost_pos++
          lost_achat_pos++
          tot_lost += -pos.diff_value
        }

        poss = append(poss, pos)
        pos = Position{}
      }
      continue
    }

    if pos.open && !pos.buy {
      if bid.Macd_signal >= -bid.Macd_absol_trigger_signal {
        pos.close_time = bid.Bid_at
        pos.close_value = bid.Last_bid
        pos.diff_value = pos.open_value - pos.close_value
        gain += pos.diff_value

        fmt.Println("Close baisse ->", pos.close_time.Format("2006-01-02 15:04:05"), "sur", pos.close_value, "with open value :", pos.open_value, "at", pos.open_time.Format("2006-01-02 15:04:05"), "gain de :", pos.diff_value)
        fmt.Println("#################################################")

        if pos.diff_value > 0 {
          win_pos++
          win_vente_pos++
          tot_win += pos.diff_value
        } else {
          lost_pos++
          lost_vente_pos++
          tot_lost += -pos.diff_value
        }

        poss = append(poss, pos)
        pos = Position{}
      }
    }
	}

  if pos.open && pos.buy {
    pos.close_time = last_bid.Bid_at
    pos.close_value = last_bid.Last_bid
    pos.diff_value = pos.close_value - pos.open_value
    gain += pos.diff_value

    fmt.Println("Close hausse ->", pos.close_time.Format("2006-01-02 15:04:05"), "on", pos.close_value, "with open value :", pos.open_value, "at", pos.open_time.Format("2006-01-02 15:04:05"), "gain de :", pos.diff_value)
    fmt.Println("#################################################")

    if pos.diff_value > 0 {
      win_pos++
      win_achat_pos++
      tot_win += pos.diff_value
    } else {
      lost_pos++
      lost_achat_pos++
      tot_lost += -pos.diff_value
    }
  } else if pos.open && !pos.buy {
    pos.close_time = last_bid.Bid_at
    pos.close_value = last_bid.Last_bid
    pos.diff_value = pos.open_value - pos.close_value
    gain += pos.diff_value

    fmt.Println("Close baisse ->", pos.close_time.Format("2006-01-02 15:04:05"), "sur", pos.close_value, "with open value :", pos.open_value, "at", pos.open_time.Format("2006-01-02 15:04:05"), "gain de :", pos.diff_value)
    fmt.Println("#################################################")

    if pos.diff_value > 0 {
      win_pos++
      win_vente_pos++
      tot_win += pos.diff_value
    } else {
      lost_pos++
      lost_vente_pos++
      tot_lost += -pos.diff_value
    }
  }


	fmt.Println("")
	fmt.Println("")
	fmt.Println("#################################################")
	fmt.Println("Total des gains final :", gain)
	fmt.Println("#######################")
	fmt.Println("Total win_pos :", win_pos, "trades")
	fmt.Println("Total win_hausse_pos :", win_achat_pos, "trades")
	fmt.Println("Total win_baisse_pos :", win_vente_pos, "trades")
	fmt.Println("Total tot_win :", tot_win, "unité")
	fmt.Println("#######################")
	fmt.Println("Total lost_pos :", lost_pos, "trades")
	fmt.Println("Total lost_hausse_pos :", lost_achat_pos, "trades")
	fmt.Println("Total lost_baisse_pos :", lost_vente_pos, "trades")
	fmt.Println("Total tot_lost :", tot_lost, "unité")
	fmt.Println("#################################################")
}
