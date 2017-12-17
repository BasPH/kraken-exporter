package main

import (
	"flag"
	"github.com/beldur/kraken-go-api-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	openingPrices = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "opening_prices",
		Help: "Opening prices",
	},
		[]string{"pair"},
	)
)

func init() {
	prometheus.MustRegister(openingPrices)
}

func main() {
	go func() {
		api := krakenapi.New("KEY", "SECRET")
		for {
			ticker, err := api.Ticker(krakenapi.XXBTZEUR)
			if err != nil {
				if ticker == nil {
					log.Warning("Result was empty. Prices not updated.")
				} else {
					log.Warning(err)
				}
			}

			askPrice, _ := strconv.ParseFloat(ticker.XXBTZEUR.Ask[0], 64)
			openingPrices.WithLabelValues("XXBTZEUR").Set(askPrice)

			time.Sleep(time.Duration(5 * time.Second))
		}
	}()

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
