package queries

import (
	"POSFlowBackend/internal/application/order/dto"
	"POSFlowBackend/internal/domain/order"
	"POSFlowBackend/internal/domain/product"
	"POSFlowBackend/internal/domain/shared"
)

type GetOrderQuery struct {
	orderRepo   order.OrderRepository
	productRepo product.ProductRepository
}

func NewGetOrderQuery(
	orderRepo order.OrderRepository,
	productRepo product.ProductRepository,
) *GetOrderQuery {
	return &GetOrderQuery{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (q *GetOrderQuery) Execute(id string) (*dto.OrderResponse, error) {
	// Find order
	ord, err := q.orderRepo.FindByID(shared.OrderID(id))
	if err != nil {
		return nil, err
	}

	// Map to response DTO
	return q.mapToDTO(ord)
}

func (q *GetOrderQuery) mapToDTO(o *order.Order) (*dto.OrderResponse, error) {
	var items []dto.OrderItemResponse

	for _, item := range o.Items() {
		prod, err := q.productRepo.FindByID(item.ProductID())
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
