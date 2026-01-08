package sqlite

import (
	"POSFlowBackend/internal/domain/sales"
	"POSFlowBackend/internal/domain/shared"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type SalesRepository struct {
	db *gorm.DB
}

func NewSalesRepository(db *gorm.DB) *SalesRepository {
	return &SalesRepository{db: db}
}

// Save implements sales.SalesRepository
func (r *SalesRepository) Save(s *sales.DailySales) error {
	model := r.toModel(s)
	return r.db.Save(&model).Error
}

// FindByID implements sales.SalesRepository
func (r *SalesRepository) FindByID(id sales.SalesID) (*sales.DailySales, error) {
	var model SalesModel

	result := r.db.Where("id = ?", id.String()).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, shared.ErrNotFound
		}
		return nil, result.Error
	}

	return r.toDomain(&model)
}

// FindByDate implements sales.SalesRepository
func (r *SalesRepository) FindByDate(date time.Time) (*sales.DailySales, error) {
	var model SalesModel

	// Normalize date to beginning of day
	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	result := r.db.Where("DATE(date) = DATE(?)", normalizedDate).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, shared.ErrNotFound
		}
		return nil, result.Error
	}

	return r.toDomain(&model)
}

// FindByDateRange implements sales.SalesRepository
func (r *SalesRepository) FindByDateRange(start, end time.Time) ([]*sales.DailySales, error) {
	var models []SalesModel

	result := r.db.Where("date BETWEEN ? AND ?", start, end).
		Order("date desc").
		Find(&models)

	if result.Error != nil {
		return nil, result.Error
	}

	return r.toDomainList(models)
}

// FindAll implements sales.SalesRepository
func (r *SalesRepository) FindAll() ([]*sales.DailySales, error) {
	var models []SalesModel

	result := r.db.Order("date desc").Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.toDomainList(models)
}

// ExistsByDate implements sales.SalesRepository
func (r *SalesRepository) ExistsByDate(date time.Time) (bool, error) {
	var count int64

	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	result := r.db.Model(&SalesModel{}).
		Where("DATE(date) = DATE(?)", normalizedDate).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

// --- Mappers: Domain Entity ↔ Database Model ---

func (r *SalesRepository) toModel(s *sales.DailySales) SalesModel {
	// Convert order IDs to JSON
	orderIDsJSON, _ := json.Marshal(s.OrderIDs())

	return SalesModel{
		ID:          s.ID().String(),
		Date:        s.Date(),
		TotalSales:  s.TotalSales().Amount,
		TotalOrders: s.TotalOrders(),
		OrderIDs:    string(orderIDsJSON),
		Closed:      s.IsClosed(),
		ClosedAt:    s.ClosedAt(),
		CreatedAt:   s.CreatedAt(),
		UpdatedAt:   s.UpdatedAt(),
	}
}

func (r *SalesRepository) toDomain(model *SalesModel) (*sales.DailySales, error) {
	// Parse order IDs from JSON
	var orderIDs []shared.OrderID
	if model.OrderIDs != "" {
		if err := json.Unmarshal([]byte(model.OrderIDs), &orderIDs); err != nil {
			return nil, err
		}
	}

	// Reconstruct domain entity with all saved values
	return sales.ReconstructDailySales(
		sales.SalesID(model.ID),
		model.Date,
		shared.Money{Amount: model.TotalSales, Currency: "USD"},
		model.TotalOrders,
		orderIDs,
		model.Closed,
		model.ClosedAt,
		model.CreatedAt,
		model.UpdatedAt,
	), nil
}

func (r *SalesRepository) toDomainList(models []SalesModel) ([]*sales.DailySales, error) {
	var result []*sales.DailySales

	for _, model := range models {
		domain, err := r.toDomain(&model)
		if err != nil {
			return nil, err
		}
		result = append(result, domain)
	}

	return result, nil
}
