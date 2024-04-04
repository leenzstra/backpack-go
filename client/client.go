package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/leenzstra/backpack-go/auth"
)

func NewBackpackClient(endpoint string, authenticator auth.Authenticator) BackpackClient {
	impl := &BackpackClientImpl{
		APIBase: APIBase{
			endpoint: endpoint,
			client:   resty.New().SetDebug(true),
		},
		Authenticator: authenticator,
	}

	impl.client.SetBaseURL(endpoint)

	impl.MarketsImpl = MarketsImpl{impl.APIBase}
	impl.SystemImpl = SystemImpl{impl.APIBase}
	impl.TradesImpl = TradesImpl{impl.APIBase}

	impl.CapitalImpl = CapitalImpl{impl.APIBase, impl.Authenticator}
	impl.HistoryImpl = HistoryImpl{impl.APIBase, impl.Authenticator}
	impl.OrdersImpl = OrdersImpl{impl.APIBase, impl.Authenticator}

	return impl
}

type APIBase struct {
	endpoint string
	client   *resty.Client
}

func (impl APIBase) Endpoint() string {
	return impl.endpoint
}

func (impl APIBase) Client() *resty.Client {
	return impl.client
}

type Base interface {
	Endpoint() string
	Client() *resty.Client
}

type BackpackClient interface {
	Markets
	System
	Trades

	Capital
	History
	Orders

	auth.Authenticator
}

type BackpackClientImpl struct {
	APIBase

	MarketsImpl
	SystemImpl
	TradesImpl

	CapitalImpl
	HistoryImpl
	OrdersImpl

	auth.Authenticator
}
