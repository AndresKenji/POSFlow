package queries

import (
	"POSFlowBackend/internal/application/product/dto"
	"POSFlowBackend/internal/domain/product"
	"POSFlowBackend/internal/domain/shared"
)

type GetProductQuery struct {
	repo product.ProductRepository
}

func NewGetProductQuery(repo product.ProductRepository) *GetProductQuery {
	return &GetProductQuery{repo: repo}
}

func (q *GetProductQuery) Execute(id string) (*dto.ProductResponse, error) {
	// Find product
	prod, err := q.repo.FindByID(shared.ProductID(id))
	if err != nil {
		return nil, err
	}

	// Map to response DTO
	return &dto.ProductResponse{
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
	}, nil
}
