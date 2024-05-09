package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

type Stock struct {
	Company string `json:"company"`
	Price   string `json:"price"`
	Change  string `json:"change"`
}

type StockData struct {
	PE   string `json:"pe"`
	EPS  string `json:"eps"`
	MCAP string `json:"mcap"`
}

type StockFull struct {
	Stock
	StockData
}

var stocksData []StockFull

func main() {
	tickers := []string{
		"MSFT",
		"IBM",
		"GE",
		"UNP",
		"COST",
		"MCD",
		"V",
		"WMT",
		"DIS",
		"MMM",
		"INTC",
		"AXP",
		"AAPL",
		"BA",
		"CSCO",
		"GS",
		"JPM",
		"CRM",
		"VZ",
	}

	c := colly.NewCollector()
	d := colly.NewCollector()

	c.OnHTML("div#quote-header-info", func(e *colly.HTMLElement) {
		stock := Stock{
			Company: e.ChildText("h1"),
			Price:   e.ChildText("fin-streamer[data-field='regularMarketPrice']"),
			Change:  e.ChildText("fin-streamer[data-field='regularMarketChangePercent']"),
		}
		stocksData = append(stocksData, StockFull{Stock: stock})
	})

	d.OnHTML("div#quote-summary", func(e *colly.HTMLElement) {
		index := len(stocksData) - 1
		if index >= 0 {
			stocksData[index].PE = e.ChildText("td[data-test='PE_RATIO-value']")
			stocksData[index].EPS = e.ChildText("td[data-test='EPS_RATIO-value']")
			stocksData[index].MCAP = e.ChildText("td[data-test='MARKET_CAP-value']")
		}
	})

	for _, ticker := range tickers {
		c.Visit("https://finance.yahoo.com/quote/" + ticker + "/")
		d.Visit("https://finance.yahoo.com/quote/" + ticker + "/")
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/stockdata", stockDataHandler)

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/home.html")
}

func stockDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(stocksData)
}
