package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"gomarket/internal/logger"
	"net/http"
	"time"
)

const (
	clientTimeout = 10

	maxAttempts     = 3
	attemptInterval = 2
)

type AccrualClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) AccrualClient {
	return AccrualClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: clientTimeout * time.Second,
		},
	}
}

func (c *AccrualClient) GetOrder(orderID string) (*OrderResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.BaseURL, orderID)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			logger.Log.Error().Err(err).Msg("err make request")

			time.Sleep(attemptInterval * time.Second)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var orderResponse OrderResponse
			err = json.NewDecoder(resp.Body).Decode(&orderResponse)
			if err != nil {
				return nil, err
			}

			if err = resp.Body.Close(); err != nil {
				return nil, fmt.Errorf("resp.Body.Close :%w", err)
			}

			return &orderResponse, nil
		}
		logger.Log.Info().Int("status_code", resp.StatusCode).Msg("unexpected status code")
		time.Sleep(attemptInterval * time.Second)
	}

	return nil, fmt.Errorf("maximum number of attempts reached, request failed")
}
