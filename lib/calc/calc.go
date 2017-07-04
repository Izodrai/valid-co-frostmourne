package calc

import (
	"../log"
	"../tools"
	"fmt"
	"time"
)

func Calc_potential(bids []tools.Bid, open_pos, close_pos *[]tools.Position) {

	fmt.Println("")
	fmt.Println("#################################################")
	log.Info("Calcul des bénéfices potentiels :")
	fmt.Println("")
	fmt.Println("")

	var last_bid tools.Bid

	for _, bid := range bids {

		last_bid = bid

		var sma_12, sma_24 float64

		for _, calc := range bid.Calculations {
			switch calc.Type {
			case "sma_12":
				sma_12 = calc.Value
			case "sma_24":
				sma_24 = calc.Value
			}
		}

		if sma_12 == 0 {
			if len(*open_pos) != 0 {
				clos_all_pos(open_pos, close_pos, bid)
				//log.Info("close all activ pos because sma_12 = ", sma_12)
			}
			continue
		}

		if sma_24 == 0 {
			if len(*open_pos) != 0 {
				clos_all_pos(open_pos, close_pos, bid)
				//log.Info("close all activ pos because sma_24 = ", sma_24)
			}
			continue
		}

		var diff_12_24 float64
		diff_12_24 = sma_24 - sma_12

		var tmp_open_pos []tools.Position

		log.Info("#######################")
		log.Info()
		log.Info("sma_24 : ", sma_24)
		log.Info("sma_12 : ", sma_12)
		log.Info("sma_24 - sma_12 : ", diff_12_24)
		log.Info()
		log.Info("bid.Bid_at : ", bid.Bid_at)
		log.Info("bid.Bid_at : ", bid.Last_bid)
		log.Info()

		/*
		for _, p := range *open_pos {

			log.Info("######")
			log.Info("p.Buy : ", p.Buy)
			log.Info("p.Open_time : ", p.Open_time)
			log.Info("p.Open_value : ", p.Open_value)

			if p.Buy {
				if diff_12_24 <= 1 {
					p.Close_time = bid.Bid_at
					p.Close_value = bid.Last_bid
					p.Diff_value = p.Close_value - p.Open_value
					p.Close_for = "diff_12_24 <= 0"

					log.Info("p.Close_time : ", p.Close_time)
					log.Info("p.Close_value : ", p.Close_value)
					log.Info("p.Diff_value : ", p.Diff_value)
					log.Info("p.Close_for : ", p.Close_for)

					*close_pos = append(*close_pos, p)
					continue
				}
			} else {
				if diff_12_24 >= -1 {
					p.Close_time = bid.Bid_at
					p.Close_value = bid.Last_bid
					p.Diff_value = p.Open_value - p.Close_value
					p.Close_for = "diff_12_24 >= 0"

					log.Info("p.Close_time : ", p.Close_time)
					log.Info("p.Close_value : ", p.Close_value)
					log.Info("p.Diff_value : ", p.Diff_value)
					log.Info("p.Close_for : ", p.Close_for)

					*close_pos = append(*close_pos, p)
					continue
				}
			}

			tmp_open_pos = append(tmp_open_pos, p)
		}

		*open_pos = tmp_open_pos

		log.Info("#######################")
		*/

		if diff_12_24 <= -1 {
			*open_pos = append(*open_pos, tools.Position{true, bid.Bid_at, bid.Last_bid, time.Time{}, 0.0, 0.0,""})
		} else if diff_12_24 >= 1 {
			*open_pos = append(*open_pos, tools.Position{false, bid.Bid_at, bid.Last_bid, time.Time{}, 0.0, 0.0,""})
		}
	}

	clos_all_pos(open_pos, close_pos, last_bid)
}

func clos_all_pos(open_pos, close_pos *[]tools.Position, bid tools.Bid) {

	for _, p := range *open_pos {

		p.Close_time = bid.Bid_at
		p.Close_value = bid.Last_bid

		if p.Buy {
			p.Diff_value = p.Close_value - p.Open_value
		} else {
			p.Diff_value = p.Open_value - p.Close_value
		}

		p.Close_for = "Close all pos"
		*close_pos = append(*close_pos, p)
	}

	*open_pos = []tools.Position{}
}

/*
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
*/
