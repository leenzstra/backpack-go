package client

import (
	"fmt"

	"github.com/leenzstra/backpack-go/auth"
)

var _ Capital = (*CapitalImpl)(nil)

type Blockchain string

const (
	TypeBitcoin  Blockchain = "Bitcoin"
	TypeEthereum Blockchain = "Ethereum"
	TypePolygon  Blockchain = "Polygon"
	TypeSolana   Blockchain = "Solana"
)

type Capital interface {
	Balances() (Balances, error)
	Deposits(limit int64, offset int64) ([]Deposit, error)
	DepositAddress(blockchain Blockchain) (*DepositAddress, error)
	Withdrawals(limit int64, offset int64) ([]Withdrawal, error)
	RequestWithdrawal(payload *WithdrawalRequest) (*Withdrawal, error)
}

type CapitalImpl struct {
	Base
	auth.Authenticator
}

// Balances implements Capital.
func (impl *CapitalImpl) Balances() (Balances, error) {
	balances := Balances{}

	headers, err := impl.Authenticate(auth.BalanceQuery, nil)
	if err != nil {
		return nil, err
	}

	resp, err := impl.Client().R().SetHeaders(headers.Map()).SetResult(&balances).Get("/api/v1/capital")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return balances, nil
}

// DepositAddress implements Capital.
func (impl *CapitalImpl) DepositAddress(blockchain Blockchain) (*DepositAddress, error) {
	deposit := &DepositAddress{}

	query := map[string]string{
		"blockchain": string(blockchain),
	}

	headers, err := impl.Authenticate(auth.DepositAddressQuery, query)
	if err != nil {
		return nil, err
	}

	resp, err := impl.Client().R().SetHeaders(headers.Map()).
		SetQueryParams(query).SetResult(deposit).Get("/wapi/v1/capital/deposit/address")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return deposit, nil
}

// Deposits implements Capital.
// Limit 0-1000
// Offset 0-N
func (impl *CapitalImpl) Deposits(limit int64, offset int64) ([]Deposit, error) {
	deposits := make([]Deposit, 0)

	query := map[string]string{
		"limit":  fmt.Sprint(limit),
		"offset": fmt.Sprint(offset),
	}

	headers, err := impl.Authenticate(auth.DepositQueryAll, query)
	if err != nil {
		return nil, err
	}

	resp, err := impl.Client().R().SetHeaders(headers.Map()).
		SetQueryParams(query).SetResult(&deposits).Get("/wapi/v1/capital/deposits")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return deposits, nil
}

// RequestWithdrawal implements Capital.
func (impl *CapitalImpl) RequestWithdrawal(payload *WithdrawalRequest) (*Withdrawal, error) {
	withdrawal := &Withdrawal{}

	headers, err := impl.Authenticate(auth.Withdraw, payload)
	if err != nil {
		return nil, err
	}

	resp, err := impl.Client().R().SetHeaders(headers.Map()).
		SetBody(payload).SetResult(withdrawal).Post("/wapi/v1/capital/withdrawals")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return withdrawal, nil
}

// Withdrawals implements Capital.
func (impl *CapitalImpl) Withdrawals(limit int64, offset int64) ([]Withdrawal, error) {
	withdrawals := make([]Withdrawal, 0)

	query := map[string]string{
		"limit":  fmt.Sprint(limit),
		"offset": fmt.Sprint(offset),
	}

	headers, err := impl.Authenticate(auth.WithdrawalQueryAll, query)
	if err != nil {
		return nil, err
	}

	resp, err := impl.Client().R().SetHeaders(headers.Map()).
		SetQueryParams(query).SetResult(&withdrawals).Get("/wapi/v1/capital/withdrawals")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return withdrawals, nil
}

type Balance struct {
	Available string `json:"available"`
	Locked    string `json:"locked"`
	Staked    string `json:"staked"`
}

type Balances map[string]Balance

type Deposit struct {
	ID                      int    `json:"id"`
	ToAddress               string `json:"toAddress"`
	FromAddress             string `json:"fromAddress"`
	ConfirmationBlockNumber int    `json:"confirmationBlockNumber"`
	ProviderID              string `json:"providerId"`
	Source                  string `json:"source"`
	Status                  string `json:"status"`
	TransactionHash         string `json:"transactionHash"`
	SubaccountID            int    `json:"subaccountId"`
	Symbol                  string `json:"symbol"`
	Quantity                string `json:"quantity"`
	CreatedAt               string `json:"createdAt"`
}

type Withdrawal struct {
	ID              int    `json:"id"`
	Blockchain      string `json:"blockchain"`
	ClientID        string `json:"clientId"`
	Identifier      string `json:"identifier"`
	Quantity        string `json:"quantity"`
	Fee             string `json:"fee"`
	Symbol          string `json:"symbol"`
	Status          string `json:"status"`
	SubaccountID    int    `json:"subaccountId"`
	ToAddress       string `json:"toAddress"`
	TransactionHash string `json:"transactionHash"`
	CreatedAt       string `json:"createdAt"`
}

type DepositAddress struct {
	Address string `json:"address"`
}

type WithdrawalRequest struct {
	Address        string `json:"address"`
	Blockchain     string `json:"blockchain"`
	ClientID       string `json:"clientId"`
	Quantity       string `json:"quantity"`
	Symbol         string `json:"symbol"`
	TwoFactorToken string `json:"twoFactorToken"`
}
