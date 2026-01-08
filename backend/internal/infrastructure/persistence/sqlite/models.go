package sqlite

import (
	"gorm.io/gorm"
	"time"
)

// ProductModel - Database representation of Product
type ProductModel struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`
	Category    string  `gorm:"not null"`
	Stock       int     `gorm:"default:0"`
	Active      bool    `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (ProductModel) TableName() string {
	return "products"
}

// OrderModel - Database representation of Order
type OrderModel struct {
	ID          string           `gorm:"primaryKey"`
	TableNumber string           `gorm:"not null"`
	Status      string           `gorm:"default:'pending'"`
	Total       float64          `gorm:"not null"`
	Items       []OrderItemModel `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (OrderModel) TableName() string {
	return "orders"
}

// OrderItemModel - Database representation of OrderItem
type OrderItemModel struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	OrderID   string  `gorm:"not null;index"`
	ProductID string  `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	UnitPrice float64 `gorm:"not null"`
	Subtotal  float64 `gorm:"not null"`
}

func (OrderItemModel) TableName() string {
	return "order_items"
}

// SalesModel - Database representation of DailySales
type SalesModel struct {
	ID          string    `gorm:"primaryKey"`
	Date        time.Time `gorm:"not null;uniqueIndex"`
	TotalSales  float64   `gorm:"default:0"`
	TotalOrders int       `gorm:"default:0"`
	OrderIDs    string    `gorm:"type:text"` // JSON array of order IDs
	Closed      bool      `gorm:"default:false"`
	ClosedAt    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (SalesModel) TableName() string {
	return "daily_sales"
}
