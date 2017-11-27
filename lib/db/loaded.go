package db

import (
	"../config"
	"../log"
	"../tools"
	"database/sql"
	"encoding/json"
	"time"
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

	rows1, err = d.Query(`SELECT COUNT(*) FROM stock_values WHERE symbol_id = ? AND calculations != "[]"`, conf.SymbolId)

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

	query := `SELECT id, symbol_id, bid_at, last_bid, calculations FROM stock_values WHERE symbol_id = ? AND calculations != "[]" ORDER BY bid_at`

	rows1, err = d.Query(query, conf.SymbolId)

	if err != nil {
		d.Close()
		return err
	}
	defer rows1.Close()

	for rows1.Next() {
		var (
			BidAt_b, Calc_b []byte
			b               tools.Bid
		)

		bar1.Increment()

		err = rows1.Scan(
			&b.Sv_id,
			&b.S_id,
			&BidAt_b,
			&b.Last_bid,
			&Calc_b)

		if err != nil {
			d.Close()
			return err
		}

		b.Bid_at, err = time.Parse("2006-01-02 15:04:05", string(BidAt_b))
		if err != nil {
			d.Close()
			return err
		}

		err = json.Unmarshal(Calc_b, &b.Calculations)
		if err != nil {
			return err
		}

		*bids = append(*bids, b)
	}

	d.Close()
	bar1.Finish()
	log.Info("Bids loaded")
	return nil
}
