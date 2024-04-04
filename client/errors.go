package client

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func extractError(req *resty.Response) error {
	return fmt.Errorf("status %d, error %s", req.StatusCode(), req.String())
}
