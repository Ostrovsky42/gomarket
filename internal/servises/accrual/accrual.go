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

type AccrualProcesser struct {
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

func NewAccrual(baseURL string, repository orders.OrderRepository) *AccrualProcesser {
	return &AccrualProcesser{
		BaseURL:         baseURL,
		OrderRepository: repository,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *AccrualProcesser) Run() {
	go a.worker()
}

func (a *AccrualProcesser) processOrder(orderID string) error {
	url := fmt.Sprintf("%s/api/orders/%s", a.BaseURL, orderID)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || resp.StatusCode == http.StatusTooManyRequests {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode) //todo more attempt
	}

	var orderResponse OrderResponse
	err = json.NewDecoder(resp.Body).Decode(&orderResponse)
	if err != nil {
		return err
	}

	errApp := a.OrderRepository.UpdateAfterAccrual(context.Background(),
		orderResponse.Order,
		orderResponse.Status,
		transferToCoins(orderResponse.Points),
	)
	if errApp != nil {
		return errApp
	}

	return nil
}

func (a *AccrualProcesser) worker() {
	for {
		orderIDs, errApp := a.OrderRepository.GetOrderIDsForAccrual(context.Background())
		if errApp != nil || len(orderIDs) == 0 {
			time.Sleep(5 * time.Second)
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

func transferToCoins(val float64) int {
	return int(100 * val)
}
