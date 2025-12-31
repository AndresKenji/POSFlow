package handlers

import (
	"POSFlowBackend/internal/application/order/commands"
	"POSFlowBackend/internal/application/order/dto"
	queries "POSFlowBackend/internal/application/order/querys"

	"POSFlowBackend/internal/infrastructure/http/request"
	"POSFlowBackend/internal/infrastructure/http/response"
	"log"

	"github.com/gin-gonic/gin"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	createCommand       *commands.CreateOrderCommand
	updateStatusCommand *commands.UpdateOrderStatusCommand
	listQuery           *queries.ListOrdersQuery
	getQuery            *queries.GetOrderQuery
	getPendingQuery     *queries.GetPendingOrdersQuery
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(
	createCommand *commands.CreateOrderCommand,
	updateStatusCommand *commands.UpdateOrderStatusCommand,
	listQuery *queries.ListOrdersQuery,
	getQuery *queries.GetOrderQuery,
	getPendingQuery *queries.GetPendingOrdersQuery,
) *OrderHandler {
	return &OrderHandler{
		createCommand:       createCommand,
		updateStatusCommand: updateStatusCommand,
		listQuery:           listQuery,
		getQuery:            getQuery,
		getPendingQuery:     getPendingQuery,
	}
}

// CreateOrder creates a new order
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest

	// Bind and validate request
	if err := request.BindAndValidate(c, &req); err != nil {
		response.HandleError(c, err)
		return
	}

	// Execute command
	order, err := h.createCommand.Execute(req)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.Created(c, order, "Order created successfully")
}

// GetOrder retrieves an order by ID
// GET /api/v1/orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := request.GetPathParam(c, "id")

	// Execute query
	order, err := h.getQuery.Execute(orderID)
	if err != nil {
		log.Printf("Error getting order: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, order, "Order retrieved successfully")
}

// ListOrders retrieves all orders
// GET /api/v1/orders
func (h *OrderHandler) ListOrders(c *gin.Context) {
	// Execute query
	orders, err := h.listQuery.Execute()
	if err != nil {
		log.Printf("Error listing orders: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, orders, "Orders retrieved successfully")
}

// GetPendingOrders retrieves all pending orders
// GET /api/v1/orders/pending
func (h *OrderHandler) GetPendingOrders(c *gin.Context) {
	// Execute query
	orders, err := h.getPendingQuery.Execute()
	if err != nil {
		log.Printf("Error getting pending orders: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, orders, "Pending orders retrieved successfully")
}

// UpdateOrderStatus updates the status of an order
// PATCH /api/v1/orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := request.GetPathParam(c, "id")

	var req dto.UpdateOrderStatusRequest

	// Bind and validate request
	if err := request.BindAndValidate(c, &req); err != nil {
		response.HandleError(c, err)
		return
	}

	// Execute command
	order, err := h.updateStatusCommand.Execute(orderID, req)
	if err != nil {
		log.Printf("Error updating order status: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, order, "Order status updated successfully")
}
