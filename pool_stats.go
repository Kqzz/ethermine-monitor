package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PoolStats struct {
	Status string `json:"status"`
	Data   struct {
		TopMiners   []interface{} `json:"topMiners"`
		MinedBlocks []struct {
			Number int    `json:"number"`
			Miner  string `json:"miner"`
			Time   int    `json:"time"`
		} `json:"minedBlocks"`
		PoolStats struct {
			HashRate      int64   `json:"hashRate"`
			Miners        int     `json:"miners"`
			Workers       int     `json:"workers"`
			BlocksPerHour float64 `json:"blocksPerHour"`
		} `json:"poolStats"`
		Price struct {
			Time time.Time `json:"time"`
			Usd  float64   `json:"usd"`
			Btc  float64   `json:"btc"`
			Eur  float64   `json:"eur"`
			Cny  float64   `json:"cny"`
			Rub  int       `json:"rub"`
		} `json:"price"`
		Estimates struct {
			Time        time.Time `json:"time"`
			BlockReward float64   `json:"blockReward"`
			Hashrate    int64     `json:"hashrate"`
			BlockTime   float64   `json:"blockTime"`
			GasPrice    float64   `json:"gasPrice"`
		} `json:"estimates"`
	} `json:"data"`
}

func GetPoolStats() (PoolStats, error) {
	var poolStats PoolStats

	res, err := http.Get("https://api.ethermine.org/poolStats")

	if err != nil {
		return poolStats, err
	}

	if res.StatusCode != 200 {
		return poolStats, fmt.Errorf("status code %d", res.StatusCode)
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&poolStats)

	if err != nil {
		return poolStats, err
	}

	return poolStats, nil
}
