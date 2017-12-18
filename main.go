package main

import (
	"github.com/beldur/kraken-go-api-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	debug = kingpin.Flag("debug", "Enable debug mode.").Short('d').Default("false").Bool()
	addr  = kingpin.Flag("address", "The address to listen on for HTTP requests").Short('a').Default(":8080").String()

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
	kingpin.Parse()

	go func() {
		key := os.Getenv("KEY")
		secret := os.Getenv("SECRET")
		api := krakenapi.New(key, secret)

		for {
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
			}

			time.Sleep(time.Duration(5 * time.Second))
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
