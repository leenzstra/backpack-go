package client

import (
	"fmt"
	"time"
)

var _ Markets = (*MarketsImpl)(nil)

type Interval string

const (
	Interval1m     Interval = "1m"
	Interval3m     Interval = "3m"
	Interval5m     Interval = "5m"
	Interval15m    Interval = "15m"
	Interval30m    Interval = "30m"
	Interval1h     Interval = "1h"
	Interval2h     Interval = "2h"
	Interval4h     Interval = "4h"
	Interval6h     Interval = "6h"
	Interval8h     Interval = "8h"
	Interval12h    Interval = "12h"
	Interval1d     Interval = "1d"
	Interval3d     Interval = "3d"
	Interval1w     Interval = "1w"
	Interval2month Interval = "1month"
)

// Public market data
type Markets interface {
	Assets() ([]Asset, error)
	Markets() ([]Market, error)
	Ticker(symbol string) (*Ticker, error)
	Tickers() ([]Ticker, error)
	Depth(symbol string) (*Depth, error)
	KLines(symbol string, interval Interval, startTime, endTime time.Time) ([]KLinePoint, error)
}

type MarketsImpl struct {
	Base
}

// Retrieves all the assets that are supported by the exchange
func (impl *MarketsImpl) Assets() ([]Asset, error) {
	assets := make([]Asset, 0)

	resp, err := impl.Client().R().SetResult(&assets).Get("/api/v1/assets")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return assets, nil
}

// GetDepth implements Markets.
func (impl *MarketsImpl) Depth(symbol string) (*Depth, error) {
	query := map[string]string{
		"symbol": symbol,
	}

	depth := &Depth{}

	resp, err := impl.Client().R().SetQueryParams(query).SetResult(&depth).Get("/api/v1/depth")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return depth, nil
}

// GetKLines implements Markets.
func (impl *MarketsImpl) KLines(symbol string, interval Interval, startTime, endTime time.Time) ([]KLinePoint, error) {
	query := map[string]string{
		"symbol":    symbol,
		"interval":  string(interval),
		"startTime": fmt.Sprint(startTime.UTC().Unix()),
		"endTime":   fmt.Sprint(endTime.UTC().Unix()),
	}

	kline := make([]KLinePoint, 0)

	resp, err := impl.Client().R().SetQueryParams(query).SetResult(&kline).Get("/api/v1/klines")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return kline, nil
}

// GetMarkets implements Markets.
func (impl *MarketsImpl) Markets() ([]Market, error) {
	markets := make([]Market, 0)

	resp, err := impl.Client().R().SetResult(&markets).Get("/api/v1/markets")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return markets, nil
}

// GetTicker implements Markets.
func (impl *MarketsImpl) Ticker(symbol string) (*Ticker, error) {
	query := map[string]string{
		"symbol": symbol,
	}

	ticker := &Ticker{}

	resp, err := impl.Client().R().SetQueryParams(query).SetResult(&ticker).Get("/api/v1/ticker")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return ticker, nil
}

// GetTickers implements Markets.
func (impl *MarketsImpl) Tickers() ([]Ticker, error) {
	tickers := make([]Ticker, 0)

	resp, err := impl.Client().R().SetResult(&tickers).Get("/api/v1/tickers")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return tickers, nil
}

type Asset struct {
	Symbol string `json:"symbol"`
	Tokens []struct {
		Blockchain        string `json:"blockchain"`
		DepositEnabled    bool   `json:"depositEnabled"`
		MinimumDeposit    string `json:"minimumDeposit"`
		WithdrawEnabled   bool   `json:"withdrawEnabled"`
		MinimumWithdrawal string `json:"minimumWithdrawal"`
		MaximumWithdrawal string `json:"maximumWithdrawal"`
		WithdrawalFee     string `json:"withdrawalFee"`
	} `json:"tokens"`
}

type Market struct {
	Symbol      string `json:"symbol"`
	BaseSymbol  string `json:"baseSymbol"`
	QuoteSymbol string `json:"quoteSymbol"`
	Filters     struct {
		Price struct {
			MinPrice string `json:"minPrice"`
			MaxPrice string `json:"maxPrice"`
			TickSize string `json:"tickSize"`
		} `json:"price"`
		Quantity struct {
			MinQuantity string `json:"minQuantity"`
			MaxQuantity string `json:"maxQuantity"`
			StepSize    string `json:"stepSize"`
		} `json:"quantity"`
		Leverage struct {
			MinLeverage string `json:"minLeverage"`
			MaxLeverage string `json:"maxLeverage"`
			StepSize    string `json:"stepSize"`
		} `json:"leverage"`
	} `json:"filters"`
}

type Ticker struct {
	Symbol             string `json:"symbol"`
	FirstPrice         string `json:"firstPrice"`
	LastPrice          string `json:"lastPrice"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	High               string `json:"high"`
	Low                string `json:"low"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	Trades             int    `json:"trades"`
}

type Depth struct {
	Asks         [][]string `json:"asks"`
	Bids         [][]string `json:"bids"`
	LastUpdateID string     `json:"lastUpdateId"`
}

type KLinePoint struct {
	Start  string `json:"start"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	End    string `json:"end"`
	Volume string `json:"volume"`
	Trades string `json:"trades"`
}
