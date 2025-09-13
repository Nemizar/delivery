package courierrepo

import (
	"github.com/google/uuid"
)

type CourierDTO struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name          string
	Speed         int
	Location      LocationDTO        `gorm:"embedded;embeddedPrefix:location_"`
	StoragePlaces []*StoragePlaceDTO `gorm:"foreignKey:CourierID;constraint:OnDelete:CASCADE"`
}

type LocationDTO struct {
	X int
	Y int
}

type StoragePlaceDTO struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	TotalVolume int
	OrderID     *uuid.UUID `gorm:"type:uuid"`
	CourierID   uuid.UUID  `gorm:"type:uuid;index"`
}

func (CourierDTO) TableName() string {
	return "couriers"
}

func (StoragePlaceDTO) TableName() string {
	return "storage_places"
}
