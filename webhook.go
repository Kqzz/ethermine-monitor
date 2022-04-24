package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"time"
)

type WebhookContent struct {
	Content     interface{}   `json:"content"`
	Embeds      []Embeds      `json:"embeds"`
	Attachments []interface{} `json:"attachments"`
}
type Author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}
type Image struct {
	URL string `json:"url"`
}
type Embeds struct {
	Color  int                      `json:"color"`
	Fields []map[string]interface{} `json:"fields"`
	Author Author                   `json:"author"`
	Image  Image                    `json:"image"`
}

func AverageHashrate(dashboard Dashboard) float64 {
	var total float64
	for _, m := range dashboard.Data.Statistics {
		total += m.CurrentHashrate
	}

	return total / float64(len(dashboard.Data.Statistics))
}

func PostWebhook(webhook string, address string, dashboard Dashboard, payouts Payouts, poolStats PoolStats, imgReader *bytes.Reader) {

	whc := WebhookContent{
		Embeds: []Embeds{
			{
				Color: 2367519,
				Author: Author{
					Name:    address,
					URL:     fmt.Sprintf("https://ethermine.org/miners/%s/dashboard", address),
					IconURL: "https://i.imgur.com/KJy2eHQ.png",
				},
				Image: Image{
					URL: "attachment://image.png",
				},
				Fields: []map[string]interface{}{
					{
						"name":   "Workers active",
						"value":  fmt.Sprintf("`%d / %d`", dashboard.Data.CurrentStatistics.ActiveWorkers, len(dashboard.Data.Workers)),
						"inline": true,
					},
					{
						"name":   "Unpaid Balance",
						"value":  fmt.Sprintf("`%f ETH` | `$%.2f USD`", WeiToEther(dashboard.Data.CurrentStatistics.Unpaid), WeiToEther(dashboard.Data.CurrentStatistics.Unpaid)*poolStats.Data.Price.Usd),
						"inline": true,
					},
					{
						"name": "Estimated Earnings",
						"value": fmt.Sprintf(
							"`%f ETH` (`$%.2f USD`) / day\n`%f ETH`  (`$%.2f USD`) / week\n`%f ETH`  (`$%.2f USD`) / month",
							payouts.Data.Estimates.CoinsPerMin*60*24,
							payouts.Data.Estimates.CoinsPerMin*60*24*poolStats.Data.Price.Usd,
							payouts.Data.Estimates.CoinsPerMin*60*24*7,
							payouts.Data.Estimates.CoinsPerMin*60*24*7*poolStats.Data.Price.Usd,
							payouts.Data.Estimates.CoinsPerMin*60*24*30,
							payouts.Data.Estimates.CoinsPerMin*60*24*30*poolStats.Data.Price.Usd,
						),
					},
					{
						"name": "————————————————————",
						"value": fmt.Sprintf(
							"**Payouts**\n\nLast Payout: `%.0f days ago`\nDaily Earnings `≈ %f  ETH`\nRemaining to Threshold: `%f  ETH`",
							math.Floor(time.Since(time.Unix(payouts.Data.Payouts[0].PaidOn, 0)).Hours()/24),
							payouts.Data.Estimates.CoinsPerMin*60*24,
							WeiToEther(dashboard.Data.Settings.MinPayout-dashboard.Data.CurrentStatistics.Unpaid),
						),
					},
					{
						"name":  "————————————————————",
						"value": "**Hashrate**",
					},
					{
						"name":   "Current",
						"value":  fmt.Sprintf("`%.1f MH/s`", HashrateToMhs(dashboard.Data.CurrentStatistics.CurrentHashrate)),
						"inline": true,
					},
					{
						"name":   "Average",
						"value":  fmt.Sprintf("`%.1f MH/s`", HashrateToMhs(AverageHashrate(dashboard))),
						"inline": true,
					},
					{
						"name":   "Reported",
						"value":  fmt.Sprintf("`%.1f MH/s`", HashrateToMhs(dashboard.Data.CurrentStatistics.ReportedHashrate)),
						"inline": true,
					},
				},
			},
		},
	}

	b, _ := json.Marshal(whc)

	req, err := newfileUploadRequest(
		webhook,
		"image",
		imgReader,
		string(b),
	)

	if err != nil {
		panic(err)
	}

	http.DefaultClient.Do(req)
}

func newfileUploadRequest(uri string, paramName string, file *bytes.Reader, jsonContent string) (*http.Request, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(paramName, "image.png")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	writer.WriteField("payload_json", jsonContent)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
