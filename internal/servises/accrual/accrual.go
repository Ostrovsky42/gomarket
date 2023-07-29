package accrual

import (
	"context"
	"gomarket/internal/external/transport"
	"gomarket/internal/logger"
	"gomarket/internal/storage/orders"
	"time"
)

const (
	workerInterval = 5
)

type Processor struct {
	transport.AccrualClient
	orders.OrderRepository
}

func NewAccrual(baseURL string, repository orders.OrderRepository) *Processor {
	return &Processor{
		OrderRepository: repository,
		AccrualClient:   transport.NewClient(baseURL),
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

		for _, orderID := range orderIDs {
			err := a.processOrder(orderID)
			if err != nil {
				logger.Log.Err(err).Str("order_id", orderID).Msg("failed process order")
			}
		}
	}
}

func (a *Processor) processOrder(orderID string) error {
	orderResponse, err := a.GetOrder(orderID)
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
