package queries

import (
	"POSFlowBackend/internal/application/product/dto"
	"POSFlowBackend/internal/domain/product"
)

type GetLowStockQuery struct {
	repo product.ProductRepository
}

func NewGetLowStockQuery(repo product.ProductRepository) *GetLowStockQuery {
	return &GetLowStockQuery{repo: repo}
}

func (q *GetLowStockQuery) Execute() (*dto.ProductListResponse, error) {
	// Find low stock products
	products, err := q.repo.FindLowStock()
	if err != nil {
		return nil, err
	}

	// Map to response DTOs
	var productResponses []*dto.ProductResponse
	for _, prod := range products {
		productResponses = append(productResponses, &dto.ProductResponse{
			ID:          prod.ID().String(),
			Name:        prod.Name(),
			Description: prod.Description(),
			Price:       prod.Price().Amount,
			Category:    string(prod.Category()),
			Stock:       prod.Stock(),
			Active:      prod.IsActive(),
			IsLowStock:  prod.IsLowStock(),
			CreatedAt:   prod.CreatedAt(),
			UpdatedAt:   prod.UpdatedAt(),
		})
	}

	return &dto.ProductListResponse{
		Products: productResponses,
		Total:    len(productResponses),
	}, nil
}
