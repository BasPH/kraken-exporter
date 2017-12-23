package main

import (
	"github.com/beldur/kraken-go-api-client"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func updateKrakenPrices(api *krakenapi.KrakenApi, intervalSec time.Duration, log *logrus.Logger) {
	for {
		log.Debug("Scraping ticker...")
		ticker, err := api.Ticker(krakenapi.XXBTZEUR)
		if err != nil {
			if ticker == nil {
				log.Warning("Result was empty. Prices not updated.")
			} else {
				log.Fatal(err)
			}
		} else {
			askPrice, _ := strconv.ParseFloat(ticker.XXBTZEUR.Ask[0], 64)
			openingPrices.WithLabelValues("XXBTZEUR").Set(askPrice)
			log.Debugf("OpeningPrice set to %v", askPrice)
		}

		time.Sleep(time.Duration(intervalSec * time.Second))
	}
}

func updateTotalBalance(api *krakenapi.KrakenApi, asset string, intervalSec time.Duration, log *logrus.Logger) {
	for {
		log.Debug("Querying Kraken...")
		result, err := api.Query("TradeBalance", map[string]string{
			"asset": asset,
		})
		if err != nil {
			log.Error(err)
		} else {
			if result != nil {
				log.Debug(result)
				r := result.(map[string]interface{})
				eb, _ := strconv.ParseFloat(r["eb"].(string), 64)
				totalBalance.WithLabelValues(asset).Set(eb)
				log.Debugf("Equivalent balance set to %v", eb)
			} else {
				log.Warning("Result was empty.")
			}
		}

		time.Sleep(time.Duration(intervalSec * time.Second))
	}
}
