package entities

import (
	"encoding/json"
	"time"
)

type Withdraw struct {
	OrderID     string    `json:"order"`
	AccountID   string    `json:"-"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (w Withdraw) MarshalJSON() ([]byte, error) {
	type WithdrawAlias Withdraw
	return json.Marshal(&struct {
		WithdrawAlias
		ProcessedAt string `json:"processed_at"`
	}{
		WithdrawAlias: WithdrawAlias(w),
		ProcessedAt:   w.ProcessedAt.Format(time.RFC3339),
	})
}
