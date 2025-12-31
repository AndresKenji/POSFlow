package product

import (
	"POSFlowBackend/internal/domain/shared"
	"time"
)

type Product struct {
	id          shared.ProductID
	name        string
	description string
	price       shared.Money
	category    Category
	stock       *Stock
	active      bool
	createdAt   time.Time
	updatedAt   time.Time
}

func NewProduct(
	id shared.ProductID,
	name string,
	price shared.Money,
	category Category,
	initialStock int,
) (*Product, error) {
	if name == "" {
		return nil, shared.ErrInvalidInput
	}

	if !category.IsValid() {
		return nil, shared.ErrInvalidInput
	}

	return &Product{
		id:        id,
		name:      name,
		price:     price,
		category:  category,
		stock:     NewStock(initialStock, 5),
		active:    true,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

func (p *Product) ID() shared.ProductID { return p.id }
func (p *Product) Name() string         { return p.name }
func (p *Product) Description() string  { return p.description }
func (p *Product) Price() shared.Money  { return p.price }
func (p *Product) Category() Category   { return p.category }
func (p *Product) Stock() int           { return p.stock.Quantity }
func (p *Product) IsActive() bool       { return p.active }
func (p *Product) CreatedAt() time.Time { return p.createdAt }
func (p *Product) UpdatedAt() time.Time { return p.updatedAt }

func (p *Product) UpdatePrice(newPrice shared.Money) error {
	if newPrice.Amount <= 0 {
		return shared.ErrInvalidPrice
	}
	p.price = newPrice
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) UpdateStock(quantity int) error {
	if quantity < 0 {
		return shared.ErrInvalidQuantity
	}

	p.stock.Quantity = quantity
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) IncreaseStock(quantity int) error {
	if quantity <= 0 {
		return shared.ErrInvalidQuantity
	}
	p.stock.Quantity += quantity
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) DecreaseStock(quantity int) error {
	if quantity <= 0 || quantity > p.stock.Quantity {
		return shared.ErrInvalidQuantity
	}
	p.stock.Quantity -= quantity
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) IsLowStock() bool {
	return p.stock.IsLowStock()
}

func (p *Product) Deactivate() {
	p.active = false
	p.updatedAt = time.Now()
}

func (p *Product) Activate() {
	p.active = true
	p.updatedAt = time.Now()
}

func (p *Product) UpdateInfo(name, description string, category Category) error {
	if name == "" {
		return shared.ErrInvalidInput
	}
	if !category.IsValid() {
		return shared.ErrInvalidInput
	}

	p.name = name
	p.description = description
	p.category = category
	p.updatedAt = time.Now()
	return nil
}
