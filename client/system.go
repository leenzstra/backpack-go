package client

import (
	"strconv"
	"time"
)

var _ System = (*SystemImpl)(nil)

type System interface {
	Status() (*Status, error)
	Ping() error
	SystemTime() (time.Time, error)
}

type SystemImpl struct {
	Base
}

// Ping implements System.
func (impl *SystemImpl) Ping() error {
	resp, err := impl.Client().R().Get("/api/v1/ping")
	if err != nil {
		return err
	}

	if resp.IsError() {
		return extractError(resp)
	}

	return nil
}

// Status implements System.
func (impl *SystemImpl) Status() (*Status, error) {
	status := &Status{}

	resp, err := impl.Client().R().SetResult(&status).Get("/api/v1/status")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, extractError(resp)
	}

	return status, nil
}

// SystemTime implements System.
func (impl *SystemImpl) SystemTime() (time.Time, error) {
	sysTime := time.Time{}

	resp, err := impl.Client().R().Get("/api/v1/time")
	if err != nil {
		return sysTime, err
	}

	if resp.IsError() {
		return sysTime, extractError(resp)
	}

	unixTime, err := strconv.ParseInt(string(resp.Body()), 10, 64)
	if err != nil {
		return sysTime, err
	}

	return time.UnixMilli(unixTime), nil
}

type Status struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
