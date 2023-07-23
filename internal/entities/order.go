package entities

import (
	"encoding/json"
	"time"
)

const (
	New        = "NEW"
	Processing = "PROCESSING"
	Invalid    = "INVALID"
	Processed  = "PROCESSED"
)

type Order struct {
	ID         string    `json:"number"`
	AccountID  string    `json:"-"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
	Points     *float64  `json:"accrual,omitempty"`
}

func (o Order) MarshalJSON() ([]byte, error) {
	type OrderAlias Order
	return json.Marshal(&struct {
		OrderAlias
		UploadedAt string `json:"uploaded_at"`
	}{
		OrderAlias: OrderAlias(o),
		UploadedAt: o.UploadedAt.Format(time.RFC3339),
	})
}
