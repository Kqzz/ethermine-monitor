package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	chart "github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
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

	if res.StatusCode != 200 {
		return dashboard, fmt.Errorf("status code %d", res.StatusCode)
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

func GetChartSeries(dashboard Dashboard) []chart.Series {
	hashrates := chart.TimeSeries{
		Name: "Current Hashrate",
		Style: chart.Style{
			StrokeColor: drawing.ColorFromHex("1f77b4"),
		},
	}

	reportedHashrates := chart.TimeSeries{
		Name: "Reported Hashrate",
		Style: chart.Style{
			StrokeColor: drawing.ColorFromHex("35c335"),
		},
	}

	for _, stat := range dashboard.Data.Statistics {
		hashrates.XValues = append(hashrates.XValues, time.Unix(int64(stat.Time), 0))
		reportedHashrates.XValues = append(reportedHashrates.XValues, time.Unix(int64(stat.Time), 0))

		hashrates.YValues = append(hashrates.YValues, HashrateToMhs(stat.CurrentHashrate))
		reportedHashrates.YValues = append(reportedHashrates.YValues, HashrateToMhs(stat.ReportedHashrate))
	}

	return []chart.Series{hashrates, reportedHashrates}
}

func GetGraph(dashboard Dashboard) *bytes.Reader {
	graph := chart.Chart{
		Width:  900,
		Height: 200,
		YAxis: chart.YAxis{
			Name: "MH/S",
			NameStyle: chart.Style{
				FontColor: drawing.ColorWhite,
			},
			Style: chart.Style{
				StrokeColor: drawing.ColorWhite,
				FontColor:   drawing.ColorWhite,
			},
		},
		XAxis: chart.XAxis{
			Name:           "Time",
			ValueFormatter: chart.TimeHourValueFormatter,
			NameStyle: chart.Style{
				FontColor: drawing.ColorWhite,
			},
			Style: chart.Style{
				StrokeColor: drawing.ColorWhite,
				FontColor:   drawing.ColorWhite,
			},
		},
		Canvas: chart.Style{
			FillColor: drawing.ColorFromHex("262327"),
			FontColor: drawing.ColorWhite,
		},
		Background: chart.Style{
			FillColor: drawing.ColorFromHex("262327"),
		},
		Series: GetChartSeries(dashboard),
	}

	buff := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buff)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(buff.Bytes())
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

	dashboard, err := GetDashboard(address)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", address)
	fmt.Printf("%.0f / %d workers running\n", dashboard.Data.CurrentStatistics.ActiveWorkers, len(dashboard.Data.Workers))
	fmt.Printf("%f ETH unpaid\n", WeiToEther(dashboard.Data.CurrentStatistics.Unpaid))
	fmt.Printf("%f MH/S\n", HashrateToMhs(dashboard.Data.CurrentStatistics.CurrentHashrate))

	body := GetGraph(dashboard)

	PostWebhook(dashboard, body)

	// time.Sleep(time.Duration(interval) * time.Minute)

}
