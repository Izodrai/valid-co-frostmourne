package calc

import (
	"../log"
	"../tools"
	"fmt"
	"time"
)

func Calc_potential(bids []tools.Bid) {

	var gain float64

	var win_pos, lost_pos int
  var win_vente_pos, lost_vente_pos int
  var win_achat_pos, lost_achat_pos int

	var tot_win, tot_lost float64

	var last_ouverture_position float64
	var last_ouverture_time time.Time

  var achat, vente bool

	fmt.Println("")
	fmt.Println("#################################################")
	log.Info("Calcul des bénéfices potentiels :")
	fmt.Println("")
	fmt.Println("")

	for _, bid := range bids {

		if last_ouverture_position == 0.0 {
			if bid.Macd_signal > bid.Macd_absol_trigger_signal {
				fmt.Println("Ouverture position à l'achat ->", bid.Bid_at.Format("2006-01-02 15:04:05"), "sur", bid.Last_bid)
				last_ouverture_position = bid.Last_bid
				last_ouverture_time = bid.Bid_at
        achat = true
			} else if bid.Macd_signal < -bid.Macd_absol_trigger_signal {
				fmt.Println("Ouverture position à la vente ->", bid.Bid_at.Format("2006-01-02 15:04:05"), "sur", bid.Last_bid)
				last_ouverture_position = bid.Last_bid
				last_ouverture_time = bid.Bid_at
        vente = true
      }
		} else {
      if achat {
        if bid.Macd_signal <= bid.Macd_absol_trigger_signal {
  				diff := bid.Last_bid - last_ouverture_position
  				fmt.Println("Fermeture position à l'achat ->", bid.Bid_at.Format("2006-01-02 15:04:05"), "sur", bid.Last_bid, "avec comme valeur d'ouverture", last_ouverture_position, "à", last_ouverture_time, "gain de :", diff)
  				fmt.Println("Total des gains jusqu'à présent :", gain)
  				gain += diff
  				fmt.Println("Après cette position :", gain)
  				fmt.Println("#################################################")
  				last_ouverture_position = 0.0
  				last_ouverture_time = time.Time{}
          achat = false
  				if diff > 0 {
  					win_pos++
            win_achat_pos++
  					tot_win += diff
  				} else {
  					lost_pos++
            lost_achat_pos++
  					tot_lost += -diff
  				}
  			}
      } else if vente {
        if bid.Macd_signal >= -bid.Macd_absol_trigger_signal {
          diff := last_ouverture_position - bid.Last_bid
  				fmt.Println("Fermeture position à la vente ->", bid.Bid_at.Format("2006-01-02 15:04:05"), "sur", bid.Last_bid, "avec comme valeur d'ouverture", last_ouverture_position, "à", last_ouverture_time, "gain de :", diff)
  				fmt.Println("Total des gains jusqu'à présent :", gain)
  				gain += diff
  				fmt.Println("Après cette position :", gain)
  				fmt.Println("#################################################")
  				last_ouverture_position = 0.0
  				last_ouverture_time = time.Time{}
          vente = false
  				if diff > 0 {
  					win_pos++
            win_vente_pos++
  					tot_win += diff
  				} else {
  					lost_pos++
            lost_vente_pos++
  					tot_lost += -diff
  				}
  			}
      }

		}

	}
	fmt.Println("")
	fmt.Println("")
	fmt.Println("#################################################")
	fmt.Println("Total des gains final :", gain)
	fmt.Println("#######################")
	fmt.Println("Total win_pos :", win_pos, "trades")
	fmt.Println("Total win_achat_pos :", win_achat_pos, "trades")
	fmt.Println("Total win_vente_pos :", win_vente_pos, "trades")
	fmt.Println("Total tot_win :", tot_win, "unité")
	fmt.Println("#######################")
	fmt.Println("Total lost_pos :", lost_pos, "trades")
	fmt.Println("Total lost_achat_pos :", lost_achat_pos, "trades")
	fmt.Println("Total lost_vente_pos :", lost_vente_pos, "trades")
	fmt.Println("Total tot_lost :", tot_lost, "unité")
	fmt.Println("#################################################")
}
