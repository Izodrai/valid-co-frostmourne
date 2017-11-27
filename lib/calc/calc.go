package calc

import (
	//"../log"
	"../tools"
	//"fmt"
	"time"
)

func Calc_potential(bids []tools.Bid, open_pos, close_pos *[]tools.Position) {

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
				clos_all_pos(open_pos, close_pos, bid, "sma_12 = 0")
			}
			continue
		}

		if sma_24 == 0 {
			if len(*open_pos) != 0 {
				clos_all_pos(open_pos, close_pos, bid, "sma_24 = 0")
			}
			continue
		}

		var diff_sma_12_24 float64
		diff_sma_12_24 = sma_24 - sma_12

		var tmp_open_pos []tools.Position

		for _, p := range *open_pos {

			if diff_sma_12_24 >= 1 {
				if !p.Buy {
					p.Close_time = bid.Bid_at
					p.Close_value = bid.Last_bid
					p.Diff_value = p.Open_value - p.Close_value
					p.Close_for = "sold pos and diff_sma_12_24 >= 1"
					p.Calculations = bid.Calculations
					*close_pos = append(*close_pos, p)
					continue
				}
			} else if diff_sma_12_24 <= -1 {
				if p.Buy {
					p.Close_time = bid.Bid_at
					p.Close_value = bid.Last_bid
					p.Diff_value = p.Close_value - p.Open_value
					p.Close_for = "buy pos and diff_sma_12_24 <= -1"
					p.Calculations = bid.Calculations
					*close_pos = append(*close_pos, p)
					continue
				}
			}

			tmp_open_pos = append(tmp_open_pos, p)
		}

		*open_pos = tmp_open_pos

		if diff_sma_12_24 >= 1 {
			*open_pos = append(*open_pos, tools.Position{true, bid.Bid_at, bid.Last_bid, time.Time{}, 0.0, 0.0, "", []tools.Calculations{}})
		} else if diff_sma_12_24 <= -1 {
			*open_pos = append(*open_pos, tools.Position{false, bid.Bid_at, bid.Last_bid, time.Time{}, 0.0, 0.0, "", []tools.Calculations{}})
		}
	}

	clos_all_pos(open_pos, close_pos, last_bid, "last_bid")
}

func clos_all_pos(open_pos, close_pos *[]tools.Position, bid tools.Bid, close_for string) {

	for _, p := range *open_pos {

		p.Close_time = bid.Bid_at
		p.Close_value = bid.Last_bid

		if p.Buy {
			p.Diff_value = p.Close_value - p.Open_value
		} else {
			p.Diff_value = p.Open_value - p.Close_value
		}

		p.Close_for = close_for
		p.Calculations = bid.Calculations

		*close_pos = append(*close_pos, p)
	}

	*open_pos = []tools.Position{}
}
