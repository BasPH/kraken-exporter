package main

import (
	"github.com/beldur/kraken-go-api-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
)

var (
	debug = kingpin.Flag("debug", "Enable debug mode.").Short('d').Default("false").Bool()
	addr  = kingpin.Flag("address", "The address to listen on for HTTP requests").Short('a').Default(":8080").String()

	api *krakenapi.KrakenApi
	log *logrus.Logger

	openingPrices = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "opening_prices",
		Help: "Opening prices",
	},
		[]string{"pair"},
	)

	totalBalance = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "total_balance",
		Help: "Total balance",
	},
		[]string{"currency"},
	)
)

func init() {
	prometheus.MustRegister(openingPrices)
	prometheus.MustRegister(totalBalance)

	log = logrus.New()

	// Init the Kraken API
	key, ok := os.LookupEnv("KEY")
	if !ok {
		log.Fatal("Environment variable KEY not set")
	}
	secret, ok := os.LookupEnv("SECRET")
	if !ok {
		log.Fatal("Environment variable SECRET not set")
	}
	api = krakenapi.New(key, secret)
}

func main() {
	kingpin.Parse()
	if *debug {
		log.SetLevel(logrus.DebugLevel)
	}
	log.Debugf("Cmd line args: %v", os.Args[1:])

	go updateTotalBalance(api, "ZEUR", 15, log)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
