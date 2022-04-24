package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Payouts struct {
	Status string `json:"status"`
	Data   struct {
		Rounds []struct {
			Block  int   `json:"block"`
			Amount int64 `json:"amount"`
		} `json:"rounds"`
		Payouts []struct {
			Start  int64       `json:"start"`
			End    int64       `json:"end"`
			Amount int64       `json:"amount"`
			TxHash string      `json:"txHash"`
			TxCost interface{} `json:"txCost"`
			PaidOn int64       `json:"paidOn"`
		} `json:"payouts"`
		PendingPayout interface{} `json:"pendingPayout"`
		MiningStart   interface{} `json:"miningStart"`
		Estimates     struct {
			AverageHashrate float64 `json:"averageHashrate"`
			CoinsPerMin     float64 `json:"coinsPerMin"`
			UsdPerMin       float64 `json:"usdPerMin"`
			BtcPerMin       float64 `json:"btcPerMin"`
		} `json:"estimates"`
	} `json:"data"`
}

func GetPayouts(address string) (Payouts, error) {
	var payouts Payouts

	res, err := http.Get(fmt.Sprintf("https://api.ethermine.org/miner/%v/dashboard/payouts", address))

	if err != nil {
		return payouts, err
	}

	if res.StatusCode != 200 {
		return payouts, fmt.Errorf("status code %d", res.StatusCode)
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&payouts)

	if err != nil {
		return payouts, err
	}

	return payouts, nil
}
