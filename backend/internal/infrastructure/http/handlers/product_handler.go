package handlers

import (
	"POSFlowBackend/internal/application/product/commands"
	"POSFlowBackend/internal/application/product/dto"
	"POSFlowBackend/internal/application/product/queries"
	"POSFlowBackend/internal/infrastructure/http/request"
	"POSFlowBackend/internal/infrastructure/http/response"
	"log"

	"github.com/gin-gonic/gin"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	createCommand      *commands.CreateProductCommand
	updateCommand      *commands.UpdateProductCommand
	deleteCommand      *commands.DeleteProductCommand
	updateStockCommand *commands.UpdateStockCommand
	listQuery          *queries.ListProductsQuery
	getQuery           *queries.GetProductQuery
	getLowStockQuery   *queries.GetLowStockQuery
}

// NewProductHandler creates a new product handler
func NewProductHandler(
	createCommand *commands.CreateProductCommand,
	updateCommand *commands.UpdateProductCommand,
	deleteCommand *commands.DeleteProductCommand,
	updateStockCommand *commands.UpdateStockCommand,
	listQuery *queries.ListProductsQuery,
	getQuery *queries.GetProductQuery,
	getLowStockQuery *queries.GetLowStockQuery,
) *ProductHandler {
	return &ProductHandler{
		createCommand:      createCommand,
		updateCommand:      updateCommand,
		deleteCommand:      deleteCommand,
		updateStockCommand: updateStockCommand,
		listQuery:          listQuery,
		getQuery:           getQuery,
		getLowStockQuery:   getLowStockQuery,
	}
}

// CreateProduct creates a new product
// POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest

	// Bind and validate request
	if err := request.BindAndValidate(c, &req); err != nil {
		response.HandleError(c, err)
		return
	}

	// Execute command
	product, err := h.createCommand.Execute(req)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.Created(c, product, "Product created successfully")
}

// GetProduct retrieves a product by ID
// GET /api/v1/products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := request.GetPathParam(c, "id")

	// Execute query
	product, err := h.getQuery.Execute(productID)
	if err != nil {
		log.Printf("Error getting product: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, product, "Product retrieved successfully")
}

// ListProducts retrieves all products
// GET /api/v1/products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	// Execute query
	products, err := h.listQuery.Execute()
	if err != nil {
		log.Printf("Error listing products: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, products, "Products retrieved successfully")
}

// UpdateProduct updates an existing product
// PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID := request.GetPathParam(c, "id")

	var req dto.UpdateProductRequest

	// Bind and validate request
	if err := request.BindAndValidate(c, &req); err != nil {
		response.HandleError(c, err)
		return
	}

	// Execute command
	product, err := h.updateCommand.Execute(productID, req)
	if err != nil {
		log.Printf("Error updating product: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, product, "Product updated successfully")
}

// DeleteProduct deletes a product (soft delete)
// DELETE /api/v1/products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID := request.GetPathParam(c, "id")

	// Execute command
	if err := h.deleteCommand.Execute(productID); err != nil {
		log.Printf("Error deleting product: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, nil, "Product deleted successfully")
}

// UpdateStock updates product stock
// POST /api/v1/products/:id/stock
func (h *ProductHandler) UpdateStock(c *gin.Context) {
	productID := request.GetPathParam(c, "id")

	var req dto.UpdateStockRequest

	// Bind and validate request
	if err := request.BindAndValidate(c, &req); err != nil {
		response.HandleError(c, err)
		return
	}

	// Execute command
	product, err := h.updateStockCommand.Execute(productID, req)
	if err != nil {
		log.Printf("Error updating stock: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, product, "Stock updated successfully")
}

// GetLowStockProducts retrieves products with low stock
// GET /api/v1/products/low-stock
func (h *ProductHandler) GetLowStockProducts(c *gin.Context) {
	// Execute query
	products, err := h.getLowStockQuery.Execute()
	if err != nil {
		log.Printf("Error getting low stock products: %v", err)
		response.HandleError(c, err)
		return
	}

	// Return success response
	response.OK(c, products, "Low stock products retrieved successfully")
}
