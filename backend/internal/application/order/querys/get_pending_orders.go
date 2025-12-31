package queries

import (
	"POSFlowBackend/internal/application/order/dto"
	"POSFlowBackend/internal/domain/order"
	"POSFlowBackend/internal/domain/product"
)

type GetPendingOrdersQuery struct {
	orderRepo   order.OrderRepository
	productRepo product.ProductRepository
}

func NewGetPendingOrdersQuery(
	orderRepo order.OrderRepository,
	productRepo product.ProductRepository,
) *GetPendingOrdersQuery {
	return &GetPendingOrdersQuery{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (q *GetPendingOrdersQuery) Execute() (*dto.OrderListResponse, error) {
	// Find pending orders
	orders, err := q.orderRepo.FindPending()
	if err != nil {
		return nil, err
	}

	// Map to response DTOs
	var orderResponses []*dto.OrderResponse
	for _, ord := range orders {
		orderDTO, err := q.mapToDTO(ord)
		if err != nil {
			return nil, err
		}
		orderResponses = append(orderResponses, orderDTO)
	}

	return &dto.OrderListResponse{
		Orders: orderResponses,
		Total:  len(orderResponses),
	}, nil
}

func (q *GetPendingOrdersQuery) mapToDTO(o *order.Order) (*dto.OrderResponse, error) {
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
