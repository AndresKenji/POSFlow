package sqlite

import (
	"POSFlowBackend/internal/domain/product"
	"POSFlowBackend/internal/domain/shared"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Save implements product.ProductRepository
func (r *ProductRepository) Save(prod *product.Product) error {
	model := r.toModel(prod)

	// Upsert: Update if exists, insert if not
	return r.db.Save(&model).Error
}

// FindByID implements product.ProductRepository
func (r *ProductRepository) FindByID(id shared.ProductID) (*product.Product, error) {
	var model ProductModel

	result := r.db.Where("id = ?", id.String()).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, shared.ErrNotFound
		}
		return nil, result.Error
	}

	return r.toDomain(&model)
}

// FindAll implements product.ProductRepository
func (r *ProductRepository) FindAll() ([]*product.Product, error) {
	var models []ProductModel

	result := r.db.Where("active = ?", true).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.toDomainList(models)
}

// FindByCategory implements product.ProductRepository
func (r *ProductRepository) FindByCategory(category product.Category) ([]*product.Product, error) {
	var models []ProductModel

	result := r.db.Where("category = ? AND active = ?", string(category), true).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.toDomainList(models)
}

// FindLowStock implements product.ProductRepository
func (r *ProductRepository) FindLowStock() ([]*product.Product, error) {
	var models []ProductModel

	// Products with stock <= 10
	result := r.db.Where("stock <= ? AND active = ?", 10, true).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.toDomainList(models)
}

// Delete implements product.ProductRepository
func (r *ProductRepository) Delete(id shared.ProductID) error {
	return r.db.Delete(&ProductModel{}, "id = ?", id.String()).Error
}

// --- Mappers: Domain Entity ↔ Database Model ---

func (r *ProductRepository) toModel(prod *product.Product) ProductModel {
	return ProductModel{
		ID:          prod.ID().String(),
		Name:        prod.Name(),
		Description: prod.Description(),
		Price:       prod.Price().Amount,
		Category:    string(prod.Category()),
		Stock:       prod.Stock(),
		Active:      prod.IsActive(),
		CreatedAt:   prod.CreatedAt(),
		UpdatedAt:   prod.UpdatedAt(),
	}
}

func (r *ProductRepository) toDomain(model *ProductModel) (*product.Product, error) {
	price, err := shared.NewMoney(model.Price)
	if err != nil {
		return nil, err
	}

	prod, err := product.NewProduct(
		shared.ProductID(model.ID),
		model.Name,
		*price,
		product.Category(model.Category),
		model.Stock,
	)
	if err != nil {
		return nil, err
	}

	// Set description if exists
	if model.Description != "" {
		prod.UpdateInfo(model.Name, model.Description, product.Category(model.Category))
	}

	// Set active status
	if !model.Active {
		prod.Deactivate()
	}

	return prod, nil
}

func (r *ProductRepository) toDomainList(models []ProductModel) ([]*product.Product, error) {
	var products []*product.Product

	for _, model := range models {
		prod, err := r.toDomain(&model)
		if err != nil {
			return nil, err
		}
		products = append(products, prod)
	}

	return products, nil
}
