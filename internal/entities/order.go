package entities

import (
	"encoding/json"
	"time"
)

const (
	New = iota + 1
	Processing
	Invalid
	Processed
)

var statusTable = map[int]string{
	New:        "NEW",
	Processing: "PROCESSING",
	Invalid:    "INVALID",
	Processed:  "PROCESSED",
}

type Order struct {
	ID         string    `json:"number"`
	AccountID  string    `json:"-"`
	Status     int       `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
	Points     *int      `json:"accrual,omitempty"`
}

func getStatus(status int) string {
	return statusTable[status]
}

func (o Order) MarshalJSON() ([]byte, error) {
	type OrderAlias Order
	return json.Marshal(&struct {
		OrderAlias
		UploadedAt string `json:"uploaded_at"`
		Status     string `json:"status"`
	}{
		OrderAlias: OrderAlias(o),
		UploadedAt: o.UploadedAt.Format(time.RFC3339),
		Status:     getStatus(o.Status),
	})
}
