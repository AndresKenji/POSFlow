package sqlite

import (
	"POSFlowBackend/internal/domain/order"
	"POSFlowBackend/internal/domain/shared"
	"time"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Save implements order.OrderRepository
func (r *OrderRepository) Save(ord *order.Order) error {
	model := r.toModel(ord)

	// Use transaction to ensure all items are saved
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing items if updating
		if err := tx.Where("order_id = ?", model.ID).Delete(&OrderItemModel{}).Error; err != nil {
			return err
		}

		// Save order with items
		return tx.Save(&model).Error
	})
}

// FindByID implements order.OrderRepository
func (r *OrderRepository) FindByID(id shared.OrderID) (*order.Order, error) {
	var model OrderModel

	result := r.db.Preload("Items").Where("id = ?", id.String()).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, shared.ErrNotFound
		}
		return nil, result.Error
	}

	return r.toDomain(&model)
}

// FindAll implements order.OrderRepository
func (r *OrderRepository) FindAll() ([]*order.Order, error) {
	var models []OrderModel

	result := r.db.Preload("Items").Order("created_at desc").Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.toDomainList(models)
}

// FindPending implements order.OrderRepository
func (r *OrderRepository) FindPending() ([]*order.Order, error) {
	var models []OrderModel

	result := r.db.Preload("Items").
		Where("status IN ?", []string{"pending", "preparing"}).
		Order("created_at asc").
		Find(&models)

	if result.Error != nil {
		return nil, result.Error
	}

	return r.toDomainList(models)
}

// FindByStatus implements order.OrderRepository
func (r *OrderRepository) FindByStatus(status order.OrderStatus) ([]*order.Order, error) {
	var models []OrderModel

	result := r.db.Preload("Items").
		Where("status = ?", string(status)).
		Order("created_at desc").
		Find(&models)

	if result.Error != nil {
		return nil, result.Error
	}

	return r.toDomainList(models)
}

// FindByDateRange implements order.OrderRepository
func (r *OrderRepository) FindByDateRange(start, end time.Time) ([]*order.Order, error) {
	var models []OrderModel

	result := r.db.Preload("Items").
		Where("created_at BETWEEN ? AND ?", start, end).
		Order("created_at desc").
		Find(&models)

	if result.Error != nil {
		return nil, result.Error
	}

	return r.toDomainList(models)
}

// --- Mappers: Domain Entity ↔ Database Model ---

func (r *OrderRepository) toModel(ord *order.Order) OrderModel {
	var items []OrderItemModel

	for _, item := range ord.Items() {
		items = append(items, OrderItemModel{
			OrderID:   ord.ID().String(),
			ProductID: item.ProductID().String(),
			Quantity:  item.Quantity(),
			UnitPrice: item.UnitPrice().Amount,
			Subtotal:  item.Subtotal().Amount,
		})
	}

	return OrderModel{
		ID:          ord.ID().String(),
		TableNumber: ord.TableNumber().String(),
		Status:      string(ord.Status()),
		Total:       ord.Total().Amount,
		Items:       items,
		CreatedAt:   ord.CreatedAt(),
		UpdatedAt:   ord.UpdatedAt(),
	}
}

func (r *OrderRepository) toDomain(model *OrderModel) (*order.Order, error) {
	var items []*order.OrderItem

	for _, itemModel := range model.Items {
		unitPrice, err := shared.NewMoney(itemModel.UnitPrice)
		if err != nil {
			return nil, err
		}

		item, err := order.NewOrderItem(
			shared.ProductID(itemModel.ProductID),
			itemModel.Quantity,
			*unitPrice,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	ord, err := order.NewOrder(
		shared.OrderID(model.ID),
		order.TableNumber(model.TableNumber),
		items,
	)
	if err != nil {
		return nil, err
	}

	// Update status if not pending
	if model.Status != "pending" {
		ord.UpdateStatus(order.OrderStatus(model.Status))
	}

	return ord, nil
}

func (r *OrderRepository) toDomainList(models []OrderModel) ([]*order.Order, error) {
	var orders []*order.Order

	for _, model := range models {
		ord, err := r.toDomain(&model)
		if err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}

	return orders, nil
}
