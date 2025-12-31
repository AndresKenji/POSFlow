package order

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPreparing OrderStatus = "preparing"
	StatusReady     OrderStatus = "ready"
	StatusCompleted OrderStatus = "completed"
	StatusCancelled OrderStatus = "cancelled"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusPreparing, StatusReady, StatusCompleted, StatusCancelled:
		return true
	}
	return false
}

func (s OrderStatus) CanTransitionTo(newStatus OrderStatus) bool {
	transitions := map[OrderStatus][]OrderStatus{
		StatusPending:   {StatusPreparing, StatusCancelled},
		StatusPreparing: {StatusReady, StatusCancelled},
		StatusReady:     {StatusCompleted, StatusCancelled},
		StatusCompleted: {},
		StatusCancelled: {},
	}

	allowedTransitions := transitions[s]
	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return true
		}
	}
	return false
}

type TableNumber string

func (t TableNumber) String() string {
	return string(t)
}
