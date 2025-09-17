package courier

import (
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrCanNotStoreOrder  = errors.New("can not store order")
	ErrStoreAlreadyClear = errors.New("store already clear")
)

type StoragePlace struct {
	baseEntity  *ddd.BaseEntity[uuid.UUID]
	name        string
	totalVolume int
	orderID     *uuid.UUID
}

func NewStoragePlace(name string, totalVolume int) (*StoragePlace, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errs.NewValueIsInvalidError("name")
	}

	if totalVolume <= 0 {
		return nil, errs.NewValueIsInvalidError("totalVolume")
	}

	return &StoragePlace{
		baseEntity:  ddd.NewBaseEntity(uuid.New()),
		name:        name,
		totalVolume: totalVolume,
	}, nil
}

func RestoreStoragePlace(id uuid.UUID, name string, totalVolume int, orderID *uuid.UUID) *StoragePlace {
	return &StoragePlace{
		baseEntity:  ddd.NewBaseEntity(id),
		name:        name,
		totalVolume: totalVolume,
		orderID:     orderID,
	}
}

func (s *StoragePlace) Id() uuid.UUID {
	return s.baseEntity.ID()
}

func (s *StoragePlace) Name() string {
	return s.name
}

func (s *StoragePlace) TotalVolume() int {
	return s.totalVolume
}

func (s *StoragePlace) OrderID() *uuid.UUID {
	return s.orderID
}

func (s *StoragePlace) Equals(other *StoragePlace) bool {
	if other == nil {
		return false
	}

	return s.baseEntity.Equal(other.baseEntity)
}

func (s *StoragePlace) CanStore(volume int) (bool, error) {
	if volume <= 0 {
		return false, errs.NewValueIsInvalidError("volume")
	}

	if s.isOccupied() {
		return false, nil
	}

	return volume <= s.totalVolume, nil
}

func (s *StoragePlace) Store(orderID uuid.UUID, volume int) error {
	if orderID == uuid.Nil {
		return errs.NewValueIsInvalidError("orderID")
	}

	canStore, err := s.CanStore(volume)
	if err != nil {
		return err
	}

	if !canStore {
		return ErrCanNotStoreOrder
	}

	s.orderID = &orderID

	return nil
}

func (s *StoragePlace) Clear() error {
	if !s.isOccupied() {
		return ErrStoreAlreadyClear
	}

	s.orderID = nil

	return nil
}

func (s *StoragePlace) isOccupied() bool {
	return s.orderID != nil
}
