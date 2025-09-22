package queries

import (
	"github.com/google/uuid"
)

type GetAllCouriersResponse struct {
	Couriers []CourierResponse
}

type CourierResponse struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name     string
	Location LocationResponse `gorm:"embedded;embeddedPrefix:location_"`
}

func (CourierResponse) TableName() string {
	return "couriers"
}
