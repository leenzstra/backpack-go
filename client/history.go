package client

import (
	"fmt"

	"github.com/leenzstra/backpack-go/auth"
)

var _ History = (*HistoryImpl)(nil)

type History interface {
	OrderHistory(orderId, symbol string, offset, limit int64) ([]Order, error)
	FillHistory(orderId, symbol string, from, to, offset, limit int64) ([]Fill, error)
}

type HistoryImpl struct {
	Base
	auth.Authenticator
}

// FillHistory implements History.
func (impl *HistoryImpl) FillHistory(orderId string, symbol string, from int64, to int64, offset int64, limit int64) ([]Fill, error) {
	history := make([]Fill, 0)

	query := map[string]string{
		"orderId": orderId,
		"symbol":  symbol,
		"from":    fmt.Sprint(from),
		"to":      fmt.Sprint(to),
		"offset":  fmt.Sprint(offset),
		"limit":   fmt.Sprint(limit),
	}

	headers, err := impl.Authenticate(auth.FillHistoryQueryAll, query)
	if err != nil {
		return nil, err
	}

	resp, err := impl.Client().R().SetHeaders(headers.Map()).SetQueryParams(query).SetResult(&history).Get("/api/v1/history/fills")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return history, nil
}

// OrderHistory implements History.
func (impl *HistoryImpl) OrderHistory(orderId string, symbol string, offset int64, limit int64) ([]Order, error) {
	history := make([]Order, 0)

	query := map[string]string{
		"orderId": orderId,
		"symbol":  symbol,
		"offset":  fmt.Sprint(offset),
		"limit":   fmt.Sprint(limit),
	}

	headers, err := impl.Authenticate(auth.OrderHistoryQueryAll, query)
	if err != nil {
		return nil, err
	}

	resp, err := impl.Client().R().SetHeaders(headers.Map()).SetQueryParams(query).SetResult(&history).Get("/api/v1/history/orders")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return history, nil
}

type Order struct {
	ID                  string `json:"id"`
	OrderType           string `json:"orderType"`
	Symbol              string `json:"symbol"`
	Side                string `json:"side"`
	Price               string `json:"price"`
	TriggerPrice        string `json:"triggerPrice"`
	Quantity            string `json:"quantity"`
	QuoteQuantity       string `json:"quoteQuantity"`
	TimeInForce         string `json:"timeInForce"`
	SelfTradePrevention string `json:"selfTradePrevention"`
	PostOnly            bool   `json:"postOnly"`
	Status              string `json:"status"`
}

type Fill struct {
	TradeID   int    `json:"tradeId"`
	OrderID   string `json:"orderId"`
	Symbol    string `json:"symbol"`
	Side      string `json:"side"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Fee       string `json:"fee"`
	FeeSymbol string `json:"feeSymbol"`
	IsMaker   bool   `json:"isMaker"`
	Timestamp string `json:"timestamp"`
}
