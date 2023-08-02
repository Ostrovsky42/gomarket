package transport

type OrderResponse struct {
	Order  string  `json:"order"`
	Status string  `json:"status"`
	Points float64 `json:"accrual"`
}
