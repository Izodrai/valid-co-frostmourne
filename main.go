package main

import (
	"./lib/calc"
	"./lib/config"
	"./lib/db"
	"./lib/log"
	"./lib/tools"
	"errors"
	"strconv"
	"fmt"
	"time"
	"strings"
	"os"
)

func init() {
	log.InitLog(false)
}

func main() {

	var err error
	var conf config.Config

	var configFile string = "config.json"

	fmt.Println("")
	fmt.Println("")

	log.YellowInfo("Running valid-co-frostmourne")

	if err = conf.LoadConfig(configFile); err != nil {
		log.FatalError(err)
		return
	}

	fmt.Println("")

	var ct_b int

	if ct_b, err = db.CountVLast(&conf); err != nil {
		log.FatalError(err)
		return
	}

	var bids []tools.Bid

	if err = db.LoadBid(&conf, ct_b, &bids); err != nil {
		log.FatalError(err)
		return
	}

	var open_pos, close_pos []tools.Position

	calc.Calc_potential(bids, &open_pos, &close_pos)

	if err = writeCSVReports(bids, close_pos); err != nil {
		log.FatalError(err)
	}
}

func writeCSVReports(bids []tools.Bid, close_pos []tools.Position) error {

	if len(bids) == 0 {
		return errors.New("No bids")
	}

	if len(close_pos) == 0 {
		return errors.New("No close_pos")
	}

	if err := os.RemoveAll("reports"); err != nil {
		return err
	}

	err := os.Mkdir("reports", 0755)
	if err != nil {
		return err
	}

	var lines_bids [][]string

	lines_bids = append(lines_bids, []string{"db_id","T Bids","V Bids","Sma_12","Sma_24","Sma_24-Sma_12"})

	for _,b := range bids {

		var sma_12, sma_24 float64

		for _, calc := range b.Calculations {
			switch calc.Type {
			case "sma_12":
				sma_12 = calc.Value
			case "sma_24":
				sma_24 = calc.Value
			}
		}

		var diff_sma_12_24 float64
		diff_sma_12_24 = sma_24 - sma_12

		if sma_12 == 0 || sma_24 == 0{
			diff_sma_12_24 = 0
		}

		lines_bids = append(lines_bids, []string{strconv.Itoa(b.Sv_id), b.Bid_at.Format("2006-01-02 15:04:05"), strconv.FormatFloat(b.Last_bid, 'f', 0, 32), strconv.FormatFloat(sma_12, 'f', 0, 32), strconv.FormatFloat(sma_24, 'f', 0, 32), strconv.FormatFloat(diff_sma_12_24, 'f', 0, 32)})
	}

	f, err := os.OpenFile("reports/bids_"+strconv.Itoa(bids[0].S_id)+"_"+time.Now().Format("2006-01-02_15:04:05")+".csv", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModeAppend|0755)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines_bids {
		if _, err := f.Write([]byte("\"" + strings.Join(line, "\";\"") + "\"\n")); err != nil {
			return err
		}
	}

	///////////////////////////////////////////////////////////////////////

	var lines_pos [][]string

	lines_pos = append(lines_pos, []string{"Position","Open Time","Close Time","Open Value","Close Value", "Diff Value", "Close For", "Sma_12", "Sma_24", "diff_sma_12_24"})

	for _,c := range close_pos {

		var pos string
		switch c.Buy {
		case true:
			pos = "buy"
		case false:
			pos = "sold"
		}

		var sma_12, sma_24 float64

		for _, calc := range c.Calculations {
			switch calc.Type {
			case "sma_12":
				sma_12 = calc.Value
			case "sma_24":
				sma_24 = calc.Value
			}
		}

		var diff_sma_12_24 float64
		diff_sma_12_24 = sma_24 - sma_12

		if sma_12 == 0 || sma_24 == 0{
			diff_sma_12_24 = 0
		}

		lines_pos = append(lines_pos, []string{pos,c.Open_time.Format("2006-01-02 15:04:05"),c.Close_time.Format("2006-01-02 15:04:05"),strconv.FormatFloat(c.Open_value, 'f', 0, 32),strconv.FormatFloat(c.Close_value, 'f', 0, 32),strconv.FormatFloat(c.Diff_value, 'f', 0, 32), c.Close_for,strconv.FormatFloat(sma_12, 'f', 0, 32),strconv.FormatFloat(sma_24, 'f', 0, 32),strconv.FormatFloat(diff_sma_12_24, 'f', 0, 32)})
	}

	f_pos, err := os.OpenFile("reports/pos_sma_"+strconv.Itoa(bids[0].S_id)+"_"+time.Now().Format("2006-01-02_15:04:05")+".csv", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModeAppend|0755)
	if err != nil {
		return err
	}
	defer f_pos.Close()

	for _, line := range lines_pos {
		if _, err := f_pos.Write([]byte("\"" + strings.Join(line, "\";\"") + "\"\n")); err != nil {
			return err
		}
	}


	return nil
}
