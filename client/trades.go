package client

import "fmt"

var _ Trades = (*TradesImpl)(nil)

type Trades interface {
	RecentTrades(symbol string, limit uint16) ([]Trade, error)
	HistoricalTrades(symbol string, limit, offset int64) ([]Trade, error)
}

type TradesImpl struct {
	Base
}

// HistoricalTrades implements Trades.
func (impl *TradesImpl) HistoricalTrades(symbol string, limit int64, offset int64) ([]Trade, error) {
	query := map[string]string{
		"symbol": symbol,
		"limit":  fmt.Sprint(limit),
		"offset": fmt.Sprint(offset),
	}

	trades := make([]Trade, 0)

	resp, err := impl.Client().R().SetQueryParams(query).SetResult(&trades).Get("/api/v1/trades/history")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return trades, nil
}

// RecentTrades implements Trades.
func (impl *TradesImpl) RecentTrades(symbol string, limit uint16) ([]Trade, error) {
	query := map[string]string{
		"symbol": symbol,
		"limit":  fmt.Sprint(limit),
	}

	trades := make([]Trade, 0)

	resp, err := impl.Client().R().SetQueryParams(query).SetResult(&trades).Get("/api/v1/trades")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return trades, nil
}

type Trade struct {
	ID            int64  `json:"id"`
	Price         string `json:"price"`
	Quantity      string `json:"quantity"`
	QuoteQuantity string `json:"quoteQuantity"`
	Timestamp     int64  `json:"timestamp"`
	IsBuyerMaker  bool   `json:"isBuyerMaker"`
}
