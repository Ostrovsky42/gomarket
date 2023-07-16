package entities

import "time"

type Order struct {
	ID         int
	AccountID  string
	Status     string
	Sum        int
	UploadedAt time.Time
}
