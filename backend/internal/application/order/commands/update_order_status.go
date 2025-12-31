package commands

import (
	"POSFlowBackend/internal/application/order/dto"
	"POSFlowBackend/internal/domain/order"
	"POSFlowBackend/internal/domain/product"
	"POSFlowBackend/internal/domain/shared"
)

type UpdateOrderStatusCommand struct {
	orderRepo   order.OrderRepository
	productRepo product.ProductRepository
}

func NewUpdateOrderStatusCommand(
	orderRepo order.OrderRepository,
	productRepo product.ProductRepository,
) *UpdateOrderStatusCommand {
	return &UpdateOrderStatusCommand{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (c *UpdateOrderStatusCommand) Execute(id string, req dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error) {
	// Find order
	ord, err := c.orderRepo.FindByID(shared.OrderID(id))
	if err != nil {
		return nil, err
	}

	// Update status using domain method
	newStatus := order.OrderStatus(req.Status)
	if err := ord.UpdateStatus(newStatus); err != nil {
		return nil, err
	}

	// Save changes
	if err := c.orderRepo.Save(ord); err != nil {
		return nil, err
	}

	// Map to response DTO
	return c.mapToDTO(ord)
}

func (c *UpdateOrderStatusCommand) mapToDTO(o *order.Order) (*dto.OrderResponse, error) {
	var items []dto.OrderItemResponse

	for _, item := range o.Items() {
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
