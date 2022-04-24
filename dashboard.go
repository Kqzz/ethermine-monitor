package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	chart "github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

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
		Width:  700,
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
			ValueFormatter: chart.TimeValueFormatterWithFormat("02 15:04"),
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

type Dashboard struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}
type Statistics struct {
	Time             int64   `json:"time"`
	LastSeen         int64   `json:"lastSeen"`
	ReportedHashrate float64 `json:"reportedHashrate"`
	CurrentHashrate  float64 `json:"currentHashrate"`
	ValidShares      int64   `json:"validShares"`
	InvalidShares    int64   `json:"invalidShares"`
	StaleShares      int64   `json:"staleShares"`
	ActiveWorkers    int64   `json:"activeWorkers"`
}
type Workers struct {
	Worker           string  `json:"worker"`
	Time             int64   `json:"time"`
	LastSeen         int64   `json:"lastSeen"`
	ReportedHashrate float64 `json:"reportedHashrate"`
	CurrentHashrate  float64 `json:"currentHashrate"`
	ValidShares      int64   `json:"validShares"`
	InvalidShares    int64   `json:"invalidShares"`
	StaleShares      int64   `json:"staleShares"`
}
type CurrentStatistics struct {
	Time             int64   `json:"time"`
	LastSeen         int64   `json:"lastSeen"`
	ReportedHashrate float64 `json:"reportedHashrate"`
	CurrentHashrate  float64 `json:"currentHashrate"`
	ValidShares      int64   `json:"validShares"`
	InvalidShares    int64   `json:"invalidShares"`
	StaleShares      int64   `json:"staleShares"`
	ActiveWorkers    int64   `json:"activeWorkers"`
	Unpaid           float64 `json:"unpaid"`
}
type Settings struct {
	Email     string  `json:"email"`
	Monitor   int64   `json:"monitor"`
	MinPayout float64 `json:"minPayout"`
	Suspended int64   `json:"suspended"`
}
type Data struct {
	Statistics        []Statistics      `json:"statistics"`
	Workers           []Workers         `json:"workers"`
	CurrentStatistics CurrentStatistics `json:"currentStatistics"`
	Settings          Settings          `json:"settings"`
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
