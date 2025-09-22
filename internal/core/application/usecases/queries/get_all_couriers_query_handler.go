package queries

import (
	"delivery/internal/pkg/errs"

	"gorm.io/gorm"
)

type GetAllCouriersQueryHandler interface {
	Handle(GetAllCouriersQuery) (GetAllCouriersResponse, error)
}

type getAllCouriersQueryHandler struct {
	db *gorm.DB
}

func NewGetAllCouriersQueryHandler(db *gorm.DB) (GetAllCouriersQueryHandler, error) {
	if db == nil {
		return &getAllCouriersQueryHandler{}, errs.NewValueIsInvalidError("db")
	}

	return &getAllCouriersQueryHandler{db: db}, nil
}

func (q *getAllCouriersQueryHandler) Handle(query GetAllCouriersQuery) (GetAllCouriersResponse, error) {
	if !query.IsValid() {
		return GetAllCouriersResponse{}, errs.NewValueIsInvalidError("query")
	}

	var couriers []CourierResponse
	result := q.db.Raw("SELECT id, name, location_x, location_y FROM couriers").Scan(&couriers)

	if result.Error != nil {
		return GetAllCouriersResponse{}, result.Error
	}

	return GetAllCouriersResponse{Couriers: couriers}, nil
}
