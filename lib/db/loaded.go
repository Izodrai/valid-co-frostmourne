package db

import (
  "time"
	"../log"
	"../config"
	"../tools"
	"database/sql"
)

func CountVLast(conf *config.Config) (int, error) {
  var ct_b int
	var err error
	var d *sql.DB
	var rows1 *sql.Rows

	if d, err = conf.DbObj.InitConnect(); err != nil {
		return 0, err
	}
	defer d.Close()

	log.WhiteInfo("Search nbr of bids in database")

	rows1, err = d.Query(`SELECT COUNT(*) FROM v_last_5_days_stock_values WHERE s_id = 1 AND sa_id IS NOT NULL`)

	if err != nil {
		return 0, err
	}
	defer rows1.Close()

	for rows1.Next() {
		err = rows1.Scan(&ct_b)
		if err != nil {
			return 0, err
		}
	}

	log.WhiteInfo("nbr of bids = ", ct_b)

	return ct_b, nil
}

func LoadBid(conf *config.Config, ct_b int, bids *[]tools.Bid) error {
  time.Sleep(1 * time.Second)

	log.Info("Load bids")

	bar1 := log.InitBar(ct_b, true)

	var err error
	var d *sql.DB
	var rows1 *sql.Rows

	if d, err = conf.DbObj.InitConnect(); err != nil {
		return err
	}

  query := "SELECT sv_id, bid_at, last_bid, sa_id, sma_c, sma_l, ema_c, ema_l, macd_value, macd_trigger, macd_signal, macd_absol_trigger_signal FROM v_last_5_days_stock_values WHERE s_id = 1 AND sa_id IS NOT NULL"

  rows1, err = d.Query(query)

	if err != nil {
		d.Close()
		return err
	}
	defer rows1.Close()

	for rows1.Next() {
		var (
			BidAt_b []byte
			b   tools.Bid
		)

    b.S_id = 1
    b.S_reference = "EURUSD"

		bar1.Increment()

		err = rows1.Scan(
			&b.Sv_id,
			&BidAt_b,
      &b.Last_bid,
      &b.Sa_id,
      &b.Sma_c,
      &b.Sma_l,
      &b.Ema_c,
      &b.Ema_l,
      &b.Macd_value,
      &b.Macd_trigger,
      &b.Macd_signal,
      &b.Macd_absol_trigger_signal)

		if err != nil {
			d.Close()
			return err
		}

		b.Bid_at, err = time.Parse("2006-01-02 15:04:05", string(BidAt_b))
		if err != nil {
			d.Close()
			return err
		}

		*bids = append(*bids, b)
	}

	d.Close()
  bar1.Finish()
	log.Info("Bids loaded")
  return nil
}
