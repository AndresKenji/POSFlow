package dto

import "time"

type DailySalesResponse struct {
	ID          string     `json:"id"`
	Date        string     `json:"date"`
	TotalSales  float64    `json:"total_sales"`
	TotalOrders int        `json:"total_orders"`
	AverageSale float64    `json:"average_sale"`
	IsClosed    bool       `json:"is_closed"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type SalesReportResponse struct {
	StartDate   string                `json:"start_date"`
	EndDate     string                `json:"end_date"`
	TotalSales  float64               `json:"total_sales"`
	TotalOrders int                   `json:"total_orders"`
	DailySales  []*DailySalesResponse `json:"daily_sales"`
	AverageSale float64               `json:"average_sale"`
}

type CloseDayResponse struct {
	Success    bool                `json:"success"`
	Message    string              `json:"message"`
	DailySales *DailySalesResponse `json:"daily_sales"`
}
