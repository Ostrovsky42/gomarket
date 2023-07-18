package entities

import "time"

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
	ID         int
	AccountID  string
	Status     int
	UploadedAt time.Time
	Points     *int
}

func GetStatus(status int) string {
	return statusTable[status]
}
