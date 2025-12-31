package commands

import (
	"POSFlowBackend/internal/application/order/dto"
	"POSFlowBackend/internal/domain/order"
	"POSFlowBackend/internal/domain/product"
	"POSFlowBackend/internal/domain/shared"

	"github.com/google/uuid"
)

type CreateOrderCommand struct {
	orderService *order.OrderService
	productRepo  product.ProductRepository
}

func NewCreateOrderCommand(
	orderService *order.OrderService,
	productRepo product.ProductRepository,
) *CreateOrderCommand {
	return &CreateOrderCommand{
		orderService: orderService,
		productRepo:  productRepo,
	}
}

func (c *CreateOrderCommand) Execute(req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	// Generate order ID
	orderID := shared.OrderID(uuid.New().String())

	// Convert DTO items to domain format
	var itemRequests []struct {
		ProductID shared.ProductID
		Quantity  int
	}

	for _, item := range req.Items {
		itemRequests = append(itemRequests, struct {
			ProductID shared.ProductID
			Quantity  int
		}{
			ProductID: shared.ProductID(item.ProductID),
			Quantity:  item.Quantity,
		})
	}

	// Use domain service to create order (handles stock validation)
	newOrder, err := c.orderService.CreateOrder(
		orderID,
		order.TableNumber(req.TableNumber),
		itemRequests,
	)
	if err != nil {
		return nil, err
	}

	// Map to response DTO
	return c.mapToDTO(newOrder)
}

func (c *CreateOrderCommand) mapToDTO(o *order.Order) (*dto.OrderResponse, error) {
	var items []dto.OrderItemResponse

	for _, item := range o.Items() {
		// Get product details
		prod, err := c.productRepo.FindByID(item.ProductID())
		if err != nil {
			return nil, err
		}

		items = append(items, dto.OrderItemResponse{
			ProductID:   item.ProductID().String(),
			ProductName: prod.Name(),
			Quantity:    item.Quantity(),
			UnitPrice:   item.UnitPrice().Amount,
			Subtotal:    item.Subtotal().Amount,
		})
	}

	return &dto.OrderResponse{
		ID:          o.ID().String(),
		TableNumber: o.TableNumber().String(),
		Status:      string(o.Status()),
		Items:       items,
		Total:       o.Total().Amount,
		CreatedAt:   o.CreatedAt(),
		UpdatedAt:   o.UpdatedAt(),
	}, nil
}
