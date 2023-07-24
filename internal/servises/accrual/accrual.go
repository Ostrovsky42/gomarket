package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"gomarket/internal/logger"
	"gomarket/internal/storage/orders"
	"net/http"
	"sync"
	"time"
)

const (
	timeout         = 10
	workerInterval  = 5
	attemptInterval = 2
)

type Processor struct {
	BaseURL    string
	HTTPClient *http.Client
	wg         sync.WaitGroup
	orders.OrderRepository
}

type OrderResponse struct {
	Order  string  `json:"order"`
	Status string  `json:"status"`
	Points float64 `json:"accrual"`
}

func NewAccrual(baseURL string, repository orders.OrderRepository) *Processor {
	return &Processor{
		BaseURL:         baseURL,
		OrderRepository: repository,
		HTTPClient: &http.Client{
			Timeout: timeout * time.Second,
		},
	}
}

func (a *Processor) Run() {
	go a.worker()
}

func (a *Processor) worker() {
	for {
		orderIDs, errApp := a.OrderRepository.GetOrderIDsForAccrual(context.Background())
		if errApp != nil || len(orderIDs) == 0 {
			time.Sleep(workerInterval * time.Second)
			continue
		}

		a.wg.Add(len(orderIDs))
		for _, orderID := range orderIDs {
			err := a.processOrder(orderID)
			if err != nil {
				logger.Log.Err(err).Str("order_id", orderID).Msg("failed process order")
			}
		}
		a.wg.Wait()
	}
}

func (a *Processor) processOrder(orderID string) error {
	url := fmt.Sprintf("%s/api/orders/%s", a.BaseURL, orderID)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	maxAttempts := 3
	for attempt := 0; attempt < maxAttempts; attempt++ {
		resp, err := a.HTTPClient.Do(req)
		if err != nil {
			logger.Log.Error().Err(err).Msg("err make request")

			time.Sleep(attemptInterval * time.Second)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var orderResponse OrderResponse
			err = json.NewDecoder(resp.Body).Decode(&orderResponse)
			if err != nil {
				return err
			}

			errApp := a.OrderRepository.UpdateAfterAccrual(context.Background(),
				orderResponse.Order,
				orderResponse.Status,
				orderResponse.Points,
			)
			if errApp != nil {
				return errApp
			}

			return nil
		}

		logger.Log.Info().Int("status_code", resp.StatusCode).Msg("unexpected status code")
		time.Sleep(attemptInterval * time.Second)
	}

	return fmt.Errorf("maximum number of attempts reached, request failed")
}
