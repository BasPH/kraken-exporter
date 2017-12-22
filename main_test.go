package main

import (
	"testing"
	"os"
	"github.com/beldur/kraken-go-api-client"
	"fmt"
)

func TestGetTradeBalanceEUR(t *testing.T) {
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	api := krakenapi.New(key, secret)

	result, err := api.Query("TradeBalance", map[string]string{
		"asset": "ZEUR",
	})
	if err != nil {
		t.Error(err)
	}

	resultMap := result.(map[string]interface{})
	fmt.Print(resultMap["eb"])
}
