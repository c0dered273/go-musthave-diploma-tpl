package clients

import (
	"net/http"
	"time"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
	"github.com/go-resty/resty/v2"
)

const (
	AccrualURL = "/api/orders/"
)

func NewAccrualClient(cfg *configs.ServerConfig) *resty.Client {
	return resty.New().
		SetBaseURL(cfg.AccrualSystemAddress).
		SetRetryCount(5).
		SetRetryMaxWaitTime(120 * time.Second).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			return response.StatusCode() == http.StatusTooManyRequests
		}).
		SetRetryAfter(func(client *resty.Client, response *resty.Response) (time.Duration, error) {
			retryAfter := response.Header().Get("Retry-After")
			duration, err := time.ParseDuration(retryAfter + "s")
			if err != nil {
				return 0, err
			}

			return duration, nil
		})
}
