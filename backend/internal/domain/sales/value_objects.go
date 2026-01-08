package sales

type SalesID string

func (s SalesID) String() string {
	return string(s)
}
