package handlers

import (
	"POSFlowBackend/internal/application/sales/commands"
	"POSFlowBackend/internal/application/sales/queries"
	"POSFlowBackend/internal/infrastructure/http/response"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// SalesHandler handles HTTP requests for sales
type SalesHandler struct {
	getDailySalesQuery  *queries.GetDailySalesQuery
	getSalesReportQuery *queries.GetSalesReportQuery
	closeDayCommand     *commands.CloseDayCommand
}

// NewSalesHandler creates a new sales handler
func NewSalesHandler(
	getDailySalesQuery *queries.GetDailySalesQuery,
	getSalesReportQuery *queries.GetSalesReportQuery,
	closeDayCommand *commands.CloseDayCommand,
) *SalesHandler {
	return &SalesHandler{
		getDailySalesQuery:  getDailySalesQuery,
		getSalesReportQuery: getSalesReportQuery,
		closeDayCommand:     closeDayCommand,
	}
}

// GetDailySales retrieves daily sales
// GET /api/v1/sales/daily?date=2026-01-08
func (h *SalesHandler) GetDailySales(c *gin.Context) {
	// Parse optional date query parameter
	var targetDate *time.Time
	dateStr := c.Query("date")

	if dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("Invalid date format: %v", err)
			response.BadRequest(c, err, "Invalid date format. Use YYYY-MM-DD")
			return
		}
		targetDate = &parsed
	}

	// Execute query
	sales, err := h.getDailySalesQuery.Execute(targetDate)
	if err != nil {
		log.Printf("Error getting daily sales: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, sales, "Daily sales retrieved successfully")
}

// GetSalesReport retrieves sales report for a date range
// GET /api/v1/sales/report?start=2026-01-01&end=2026-01-31
func (h *SalesHandler) GetSalesReport(c *gin.Context) {
	// Parse required date parameters
	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		response.BadRequest(c, errors.New("missing required parameters"), "start and end dates are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		log.Printf("Invalid start date: %v", err)
		response.BadRequest(c, err, "Invalid start date format. Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		log.Printf("Invalid end date: %v", err)
		response.BadRequest(c, err, "Invalid end date format. Use YYYY-MM-DD")
		return
	}

	// Execute query
	report, err := h.getSalesReportQuery.Execute(startDate, endDate)
	if err != nil {
		log.Printf("Error getting sales report: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, report, "Sales report retrieved successfully")
}

// CloseDay closes the day and generates final report
// POST /api/v1/sales/close-day
func (h *SalesHandler) CloseDay(c *gin.Context) {
	// Optional: allow closing a specific date via query parameter
	var targetDate *time.Time
	dateStr := c.Query("date")

	if dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("Invalid date format: %v", err)
			response.BadRequest(c, err, "Invalid date format. Use YYYY-MM-DD")
			return
		}
		targetDate = &parsed
	}

	// Execute command
	result, err := h.closeDayCommand.Execute(targetDate)
	if err != nil {
		log.Printf("Error closing day: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, result, result.Message)
}
