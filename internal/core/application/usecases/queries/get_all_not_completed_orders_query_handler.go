package queries

import (
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"

	"gorm.io/gorm"
)

type GetNotCompletedOrdersQueryHandler interface {
	Handle(GetNotCompletedOrdersQuery) (GetNotCompletedOrdersResponse, error)
}

type getNotCompletedOrdersQueryHandler struct {
	db *gorm.DB
}

func NewGetNotCompletedOrdersQueryHandler(db *gorm.DB) (GetNotCompletedOrdersQueryHandler, error) {
	if db == nil {
		return &getNotCompletedOrdersQueryHandler{}, errs.NewValueIsInvalidError("db")
	}
	return &getNotCompletedOrdersQueryHandler{db: db}, nil
}

func (q *getNotCompletedOrdersQueryHandler) Handle(query GetNotCompletedOrdersQuery) (GetNotCompletedOrdersResponse, error) {
	if !query.IsValid() {
		return GetNotCompletedOrdersResponse{}, errs.NewValueIsInvalidError("query")
	}

	var orders []OrderResponse
	result := q.db.Raw("SELECT id, location_x, location_y FROM orders WHERE status != ?", order.StatusCompleted).Scan(&orders)

	if result.Error != nil {
		return GetNotCompletedOrdersResponse{}, result.Error
	}

	return GetNotCompletedOrdersResponse{Orders: orders}, nil
}
