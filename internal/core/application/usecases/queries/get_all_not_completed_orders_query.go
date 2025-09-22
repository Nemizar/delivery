package queries

type GetNotCompletedOrdersQuery struct {
	isValid bool
}

func NewGetNotCompletedOrdersQuery() (GetNotCompletedOrdersQuery, error) {
	return GetNotCompletedOrdersQuery{isValid: true}, nil
}

func (q GetNotCompletedOrdersQuery) IsValid() bool {
	return q.isValid
}
