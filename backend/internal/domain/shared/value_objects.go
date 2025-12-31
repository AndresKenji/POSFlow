package shared

import "fmt"

type Money struct {
	Amount   float64
	Currency string
}

func NewMoney(amount float64) (*Money, error) {
	if amount < 0 {
		return nil, ErrInvalidPrice
	}
	return &Money{
		Amount:   amount,
		Currency: "USD",
	}, nil
}

func (m *Money) Add(other Money) Money {
	return Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}
}

func (m *Money) Multiply(factor float64) Money {
	return Money{
		Amount:   m.Amount * factor,
		Currency: m.Currency,
	}
}

func (m *Money) String() string {
	return fmt.Sprintf("%.2f %s", m.Amount, m.Currency)
}

type ProductID string

func (p ProductID) String() string {
	return string(p)
}

type OrderID string

func (o OrderID) String() string {
	return string(o)
}
