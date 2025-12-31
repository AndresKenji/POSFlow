package queries

import (
	"POSFlowBackend/internal/application/product/dto"
	"POSFlowBackend/internal/domain/product"
)

type ListProductsQuery struct {
	repo product.ProductRepository
}

func NewListProductsQuery(repo product.ProductRepository) *ListProductsQuery {
	return &ListProductsQuery{repo: repo}
}

func (q *ListProductsQuery) Execute() (*dto.ProductListResponse, error) {
	// Find all products
	products, err := q.repo.FindAll()
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
