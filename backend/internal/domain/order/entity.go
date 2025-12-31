package order

import (
	"POSFlowBackend/internal/domain/shared"
	"time"
)

// OrderItem represents an item in an order
type OrderItem struct {
	productID shared.ProductID
	quantity  int
	unitPrice shared.Money
	subtotal  shared.Money
}

func NewOrderItem(productID shared.ProductID, quantity int, unitPrice shared.Money) (*OrderItem, error) {
	if quantity <= 0 {
		return nil, shared.ErrInvalidQuantity
	}

	subtotal := unitPrice.Multiply(float64(quantity))

	return &OrderItem{
		productID: productID,
		quantity:  quantity,
		unitPrice: unitPrice,
		subtotal:  subtotal,
	}, nil
}

func (oi *OrderItem) ProductID() shared.ProductID { return oi.productID }
func (oi *OrderItem) Quantity() int               { return oi.quantity }
func (oi *OrderItem) UnitPrice() shared.Money     { return oi.unitPrice }
func (oi *OrderItem) Subtotal() shared.Money      { return oi.subtotal }

// Order is an aggregate root
type Order struct {
	id          shared.OrderID
	tableNumber TableNumber
	items       []*OrderItem
	status      OrderStatus
	total       shared.Money
	createdAt   time.Time
	updatedAt   time.Time
}

// NewOrder creates a new Order (Factory method)
func NewOrder(
	id shared.OrderID,
	tableNumber TableNumber,
	items []*OrderItem,
) (*Order, error) {

	if len(items) == 0 {
		return nil, shared.ErrInvalidInput
	}

	// Calculate total
	total := shared.Money{Amount: 0, Currency: "USD"}
	for _, item := range items {
		total = total.Add(item.subtotal)
	}

	return &Order{
		id:          id,
		tableNumber: tableNumber,
		items:       items,
		status:      StatusPending,
		total:       total,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}, nil
}

// Getters
func (o *Order) ID() shared.OrderID       { return o.id }
func (o *Order) TableNumber() TableNumber { return o.tableNumber }
func (o *Order) Items() []*OrderItem      { return o.items }
func (o *Order) Status() OrderStatus      { return o.status }
func (o *Order) Total() shared.Money      { return o.total }
func (o *Order) CreatedAt() time.Time     { return o.createdAt }
func (o *Order) UpdatedAt() time.Time     { return o.updatedAt }

// Business methods
func (o *Order) UpdateStatus(newStatus OrderStatus) error {
	if !newStatus.IsValid() {
		return shared.ErrInvalidInput
	}

	if !o.status.CanTransitionTo(newStatus) {
		return shared.ErrOrderNotModifiable
	}

	o.status = newStatus
	o.updatedAt = time.Now()
	return nil
}

func (o *Order) StartPreparing() error {
	return o.UpdateStatus(StatusPreparing)
}

func (o *Order) MarkReady() error {
	return o.UpdateStatus(StatusReady)
}

func (o *Order) Complete() error {
	return o.UpdateStatus(StatusCompleted)
}

func (o *Order) Cancel() error {
	return o.UpdateStatus(StatusCancelled)
}

func (o *Order) IsPending() bool {
	return o.status == StatusPending || o.status == StatusPreparing
}

func (o *Order) IsCompleted() bool {
	return o.status == StatusCompleted
}
