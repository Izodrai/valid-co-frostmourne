package main

import (
	"./lib/calc"
	"./lib/config"
	"./lib/db"
	"./lib/log"
	"./lib/tools"
	"fmt"
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

	//var total float64
	/*
	for _, close_pos := range close_pos {
		log.Info(close_pos)
	}
	*/
}
