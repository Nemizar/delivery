package queries

type GetAllCouriersQuery struct {
	isValid bool
}

func NewGetAllCouriersQuery() (GetAllCouriersQuery, error) {
	return GetAllCouriersQuery{isValid: true}, nil
}

func (q GetAllCouriersQuery) IsValid() bool {
	return q.isValid
}
