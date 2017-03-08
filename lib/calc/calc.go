package calc

import (
  "../tools"
	"../log"
  "fmt"
  "time"
)

func Calc_potential (bids []tools.Bid) {

  var gain float64

  var win_pos, lost_pos int
  var tot_win, tot_lost float64

  var last_ouverture_position float64
  var last_ouverture_time time.Time

  fmt.Println("")
  fmt.Println("#################################################")
  log.Info("Calcul des bénéfices potentiels :")
  fmt.Println("")
  fmt.Println("")

  for _, bid := range bids {

    if last_ouverture_position == 0.0 {
      if bid.Macd_signal > bid.Macd_absol_trigger_signal {
        fmt.Println("Ouverture position ->", bid.Bid_at.Format("2006-01-02 15:04:05"), "sur", bid.Last_bid)
        last_ouverture_position = bid.Last_bid
        last_ouverture_time = bid.Bid_at
      }
    } else {
      if bid.Macd_signal <= bid.Macd_absol_trigger_signal {
        diff := bid.Last_bid - last_ouverture_position
        fmt.Println("Fermeture position ->", bid.Bid_at.Format("2006-01-02 15:04:05"), "sur", bid.Last_bid, "avec comme valeur d'ouverture", last_ouverture_position, "à", last_ouverture_time, "gain de :", diff)
        fmt.Println("Total des gains jusqu'à présent :", gain)
        gain += diff
        fmt.Println("Après cette position :", gain)
        fmt.Println("#################################################")
        last_ouverture_position = 0.0
        last_ouverture_time = time.Time{}
        if diff > 0 {
          win_pos++
          tot_win += diff
        } else {
          lost_pos++
          tot_lost += -diff
        }
      }
    }


    /*
    if bid.Macd_signal > bid.Macd_absol_trigger_signal || bid.Macd_signal < - bid.Macd_absol_trigger_signal {

      if
      fmt.Println(bid.Sv_id, " - ", bid.Bid_at, " - ", bid.Macd_signal, " - ", bid.Macd_absol_trigger_signal )
    }
    */

  }
  fmt.Println("")
  fmt.Println("")
  fmt.Println("#################################################")
  fmt.Println("Total des gains final :", gain)
  fmt.Println("Total win_pos :", win_pos)
  fmt.Println("Total tot_win :", tot_win)
  fmt.Println("#######################")
  fmt.Println("Total lost_pos :", lost_pos)
  fmt.Println("Total tot_lost :", tot_lost)
  fmt.Println("#################################################")
}
