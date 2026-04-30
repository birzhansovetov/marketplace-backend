package client

import (
	"time"

	"github.com/go-resty/resty/v2"
)

func NewRestyClient() *resty.Client {
	return resty.New().
		SetTimeout(10*time.Second).
		SetHeader("Content-Type", "application/json").
		SetRetryCount(3).
		SetRetryWaitTime(500 * time.Millisecond)
}
