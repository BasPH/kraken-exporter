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
	pairs = kingpin.Flag("pairs", "Pairs to fetch from Kraken").Short('p').Default("XXBTZEUR").String()

	openingPrices = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "opening_prices",
		Help: "Opening prices",
	},
		[]string{"pair"},
	)

	totalBalance = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "total_balance",
		Help: "Total balance",
	})
)

func init() {
	prometheus.MustRegister(openingPrices)
	prometheus.MustRegister(totalBalance)
}

func fetchKrakenPrices() {
	key, ok := os.LookupEnv("KEY")
	if !ok {
		log.Fatal("Environment variable KEY not set")
	}
	secret, ok := os.LookupEnv("SECRET")
	if !ok {
		log.Fatal("Environment variable SECRET not set")
	}
	api := krakenapi.New(key, secret)

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

		time.Sleep(time.Duration(5 * time.Second))
	}
}

func fetchTotalBalance() {
	key, ok := os.LookupEnv("KEY")
	if !ok {
		log.Fatal("Environment variable KEY not set")
	}
	secret, ok := os.LookupEnv("SECRET")
	if !ok {
		log.Fatal("Environment variable SECRET not set")
	}
	api := krakenapi.New(key, secret)

	for {
		log.Debug("Querying Kraken...")
		result, err := api.Query("TradeBalance", map[string]string{
			"asset": "ZEUR",
		})
		if err != nil {
			log.Error(err)
		} else {
			if result != nil {
				r := result.(map[string]interface{})
				eb := r["eb"].(float64)
				totalBalance.Set(eb)
				log.Debugf("Equivalent balance set to %v", eb)
			} else {
				log.Warning("Result was empty.")
			}
		}

		time.Sleep(time.Duration(5 * time.Second))
	}
}

func main() {
	kingpin.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("Cmd line args: %v", os.Args[1:])

	//go fetchKrakenPrices()
	go fetchTotalBalance()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
