package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Dashboard struct {
	Status string `json:"status"`
	Data   struct {
		Statistics []struct {
			Time             float64 `json:"time"`
			LastSeen         float64 `json:"lastSeen"`
			ReportedHashrate float64 `json:"reportedHashrate"`
			CurrentHashrate  float64 `json:"currentHashrate"`
			ValidShares      float64 `json:"validShares"`
			InvalidShares    float64 `json:"invalidShares"`
			StaleShares      float64 `json:"staleShares"`
			ActiveWorkers    float64 `json:"activeWorkers"`
		} `json:"statistics"`
		Workers []struct {
			Worker           string  `json:"worker"`
			Time             float64 `json:"time"`
			LastSeen         float64 `json:"lastSeen"`
			ReportedHashrate float64 `json:"reportedHashrate"`
			CurrentHashrate  float64 `json:"currentHashrate"`
			ValidShares      float64 `json:"validShares"`
			InvalidShares    float64 `json:"invalidShares"`
			StaleShares      float64 `json:"staleShares"`
		} `json:"workers"`
		CurrentStatistics struct {
			Time             float64 `json:"time"`
			LastSeen         float64 `json:"lastSeen"`
			ReportedHashrate float64 `json:"reportedHashrate"`
			CurrentHashrate  float64 `json:"currentHashrate"`
			ValidShares      float64 `json:"validShares"`
			InvalidShares    float64 `json:"invalidShares"`
			StaleShares      float64 `json:"staleShares"`
			ActiveWorkers    float64 `json:"activeWorkers"`
			Unpaid           float64 `json:"unpaid"`
		} `json:"currentStatistics"`
		Settings struct {
			Email     string  `json:"email"`
			Monitor   float64 `json:"monitor"`
			MinPayout float64 `json:"minPayout"`
			Suspended float64 `json:"suspended"`
		} `json:"settings"`
	} `json:"data"`
}

func GetDashboard(address string) (Dashboard, error) {
	var dashboard Dashboard

	res, err := http.Get(fmt.Sprintf("https://api.ethermine.org/miner/%v/dashboard", address))

	if err != nil {
		return dashboard, err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&dashboard)

	if err != nil {
		return dashboard, err
	}

	return dashboard, nil
}

func WeiToEther(wei float64) float64 {
	return wei / 1000000000000000000
}

func HashrateToMhs(hashrate float64) float64 {
	return hashrate / 1000000
}

func main() {
	address := "0xf87ca7d6b9ac4f656fa73b180bfe3c42b38eede1"
	dashboard, err := GetDashboard(address)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", address)
	fmt.Printf("%.0f / %d workers running\n", dashboard.Data.CurrentStatistics.ActiveWorkers, len(dashboard.Data.Workers))
	fmt.Printf("%f ETH unpaid\n", WeiToEther(dashboard.Data.CurrentStatistics.Unpaid))
	fmt.Printf("%f MH/S\n", HashrateToMhs(dashboard.Data.CurrentStatistics.CurrentHashrate))

}
