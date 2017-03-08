package main

import(
  "fmt"
	"./lib/db"
	"./lib/log"
	"./lib/calc"
	"./lib/tools"
	"./lib/config"
)

func init() {
	log.InitLog(false)
}

func main() {
  var err error
	var conf config.Config

	var configFile string = "config/conf.json"

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

  calc.Calc_potential(bids)
}
