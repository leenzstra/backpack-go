package client

import (
	"fmt"

	"github.com/leenzstra/backpack-go/auth"
)

var _ Orders = (*OrdersImpl)(nil)

type Orders interface {
	OpenOrder(clientId uint32, orderId string, symbol string) (*BaseOrder, error)
	ExecuteOrder(payload ExecuteOrderPayload) (*BaseOrder, error)
	CancelOrder(payload CancelOrderPayload) (*BaseOrder, error)

	OpenOrders(symbol string) ([]BaseOrder, error)
	CancelOrders(payload CancelOrderPayload) ([]BaseOrder, error)
}

type OrdersImpl struct {
	Base
	auth.Authenticator
}

// CancelOrder implements Orders.
func (impl *OrdersImpl) CancelOrder(payload CancelOrderPayload) (*BaseOrder, error) {
	order := &BaseOrder{}

	headers, err := impl.Authenticate(auth.OrderCancel, payload)
	if err != nil {
		return nil, err
	}

	resp, err := impl.Client().R().SetHeaders(headers.Map()).
		SetBody(payload).SetResult(order).Delete("/api/v1/order")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return order, nil
}

// CancelOrders implements Orders.
//
// Fill only symbol
func (impl *OrdersImpl) CancelOrders(payload CancelOrderPayload) ([]BaseOrder, error) {
	orders := make([]BaseOrder, 0)

	required := map[string]string{
		"symbol": payload.Symbol,
	}

	headers, err := impl.Authenticate(auth.OrderCancelAll, required)
	if err != nil {
		return nil, err
	}

	resp, err := impl.Client().R().SetHeaders(headers.Map()).
		SetBody(required).SetResult(&orders).Delete("/api/v1/orders")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return orders, nil
}

// ExecuteOrder implements Orders.
func (impl *OrdersImpl) ExecuteOrder(payload ExecuteOrderPayload) (*BaseOrder, error) {
	order := &BaseOrder{}

	headers, err := impl.Authenticate(auth.OrderExecute, payload)
	if err != nil {
		return nil, err
	}

	resp, err := impl.Client().R().SetHeaders(headers.Map()).
		SetBody(payload).SetResult(order).Post("/api/v1/order")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return order, nil
}

// OpenOrder implements Orders.
func (impl *OrdersImpl) OpenOrder(clientId uint32, orderId string, symbol string) (*BaseOrder, error) {
	order := &BaseOrder{}

	query := map[string]string{
		"clientId": fmt.Sprint(clientId),
		"orderId":  orderId,
		"symbol":   symbol,
	}

	headers, err := impl.Authenticate(auth.OrderQuery, query)
	if err != nil {
		return nil, err
	}

	resp, err := impl.Client().R().SetHeaders(headers.Map()).
		SetQueryParams(query).SetResult(order).Get("/api/v1/order")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return order, nil
}

// OpenOrders implements Orders.
func (impl *OrdersImpl) OpenOrders(symbol string) ([]BaseOrder, error) {
	orders := make([]BaseOrder, 0)

	query := map[string]string{
		"symbol": symbol,
	}

	headers, err := impl.Authenticate(auth.OrderQueryAll, query)
	if err != nil {
		return nil, err
	}

	resp, err := impl.Client().R().SetHeaders(headers.Map()).
		SetQueryParams(query).SetResult(&orders).Get("/api/v1/orders")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return orders, nil
}

type BaseOrder struct {
	OrderType             string `json:"orderType"`
	ID                    string `json:"id"`
	ClientID              int    `json:"clientId"`
	Symbol                string `json:"symbol"`
	Side                  string `json:"side"`
	Quantity              string `json:"quantity"`
	ExecutedQuantity      string `json:"executedQuantity"`
	ExecutedQuoteQuantity string `json:"executedQuoteQuantity"`
	TriggerPrice          string `json:"triggerPrice"`
	TimeInForce           string `json:"timeInForce"`
	SelfTradePrevention   string `json:"selfTradePrevention"`
	Status                string `json:"status"`
	CreatedAt             int    `json:"createdAt"`
}

type MarketOrder struct {
	BaseOrder
	QuoteQuantity string `json:"quoteQuantity"`
}

type LimitOrder struct {
	BaseOrder
	Price    string `json:"price"`
	PostOnly bool   `json:"postOnly"`
}

type ExecuteOrderPayload struct {
	ClientID            int    `json:"clientId"`
	OrderType           string `json:"orderType"`
	PostOnly            bool   `json:"postOnly"`
	Price               string `json:"price"`
	Quantity            string `json:"quantity"`
	QuoteQuantity       string `json:"quoteQuantity"`
	SelfTradePrevention string `json:"selfTradePrevention"`
	Side                string `json:"side"`
	Symbol              string `json:"symbol"`
	TimeInForce         string `json:"timeInForce"`
	TriggerPrice        string `json:"triggerPrice"`
}

type CancelOrderPayload struct {
	ClientID int    `json:"clientId"`
	OrderID  string `json:"orderId"`
	Symbol   string `json:"symbol"`
}
