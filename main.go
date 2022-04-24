package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func WeiToEther(wei float64) float64 {
	return wei / 1000000000000000000
}

func HashrateToMhs(hashrate float64) float64 {
	return hashrate / 1000000
}

func main() {

	var interval float64

	flag.Float64Var(&interval, "interval", 5, "Interval in minutes")

	flag.Parse()

	address := flag.Arg(0)

	if address == "" {
		fmt.Println("no address given")
		os.Exit(0)
	}

	for {

		dashboard, err := GetDashboard(address)

		if err != nil {
			panic(err)
		}

		poolStats, err := GetPoolStats()

		if err != nil {
			panic(err)
		}

		payouts, err := GetPayouts(address)

		if err != nil {
			panic(err)
		}

		fmt.Printf("%v\n", address)
		fmt.Printf("%d / %d workers running\n", dashboard.Data.CurrentStatistics.ActiveWorkers, len(dashboard.Data.Workers))
		fmt.Printf("%f ETH unpaid\n", WeiToEther(dashboard.Data.CurrentStatistics.Unpaid))
		fmt.Printf("%f MH/S\n", HashrateToMhs(dashboard.Data.CurrentStatistics.CurrentHashrate))

		body := GetGraph(dashboard)

		PostWebhook(address, dashboard, payouts, poolStats, body)

		time.Sleep(time.Duration(interval) * time.Minute)

	}

}
