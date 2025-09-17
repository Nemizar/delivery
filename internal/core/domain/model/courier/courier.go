package courier

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"errors"
	"math"

	"github.com/google/uuid"
)

var (
	ErrNoSuitablePlace = errors.New("no suitable place")
	ErrOrderNotFound   = errors.New("order not found")
)

type Courier struct {
	baseAggregate *ddd.BaseAggregate[uuid.UUID]
	name          string
	speed         int
	location      kernel.Location
	storagePlaces []*StoragePlace
}

func NewCourier(name string, speed int, location kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, errs.NewValueIsInvalidError("name")
	}

	if speed <= 0 {
		return nil, errs.NewValueIsInvalidError("speed")
	}

	if !location.IsValid() {
		return nil, errs.NewValueIsInvalidError("location")
	}

	c := &Courier{
		baseAggregate: ddd.NewBaseAggregate(uuid.New()),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: make([]*StoragePlace, 0),
	}

	err := c.AddStoragePlace("Bag", 10)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func RestoreCourier(id uuid.UUID, name string, speed int, location kernel.Location, storagePlaces []*StoragePlace) *Courier {
	return &Courier{
		baseAggregate: ddd.NewBaseAggregate(id),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: storagePlaces,
	}
}

func (c *Courier) Id() uuid.UUID {
	return c.baseAggregate.ID()
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Speed() int {
	return c.speed
}

func (c *Courier) Location() kernel.Location {
	return c.location
}

func (c *Courier) StoragePlaces() []*StoragePlace {
	return c.storagePlaces
}

func (c *Courier) Equals(other *Courier) bool {
	if other == nil {
		return false
	}

	return c.baseAggregate.Equal(other.baseAggregate)
}

func (c *Courier) AddStoragePlace(name string, volume int) error {
	sp, err := NewStoragePlace(name, volume)
	if err != nil {
		return err
	}

	c.storagePlaces = append(c.storagePlaces, sp)

	return nil
}

func (c *Courier) CanTakeOrder(order *order.Order) (bool, error) {
	if order == nil {
		return false, errs.NewValueIsInvalidError("order")
	}

	for _, place := range c.storagePlaces {
		if place.orderID != nil {
			continue
		}

		can, err := place.CanStore(order.Volume())
		if err != nil {
			return false, err
		}

		return can, nil
	}

	return false, nil
}

func (c *Courier) TakeOrder(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsInvalidError("order")
	}

	for _, place := range c.storagePlaces {
		if place.orderID != nil {
			continue
		}

		can, err := place.CanStore(order.Volume())
		if err != nil {
			return err
		}

		if !can {
			return ErrNoSuitablePlace
		}

		err = place.Store(order.ID(), order.Volume())
		if err != nil {
			return err
		}

		return nil
	}

	return ErrNoSuitablePlace
}

func (c *Courier) CompleteOrder(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsInvalidError("order")
	}

	sp, err := c.findStoragePlaceByOrderID(order.ID())
	if err != nil {
		return err
	}

	err = sp.Clear()
	if err != nil {
		return err
	}

	return nil
}

func (c *Courier) CalculateTimeToLocation(target kernel.Location) (float64, error) {
	if !target.IsValid() {
		return 0, errs.NewValueIsInvalidError("target")
	}

	d := c.location.DistanceTo(target)

	t := float64(d) / float64(c.speed)

	return t, nil
}

func (c *Courier) Move(target kernel.Location) error {
	if !target.IsValid() {
		return errs.NewValueIsRequiredError("target")
	}

	dx := float64(target.X() - c.location.X())
	dy := float64(target.Y() - c.location.Y())
	remainingRange := float64(c.speed)

	if math.Abs(dx) > remainingRange {
		dx = math.Copysign(remainingRange, dx)
	}
	remainingRange -= math.Abs(dx)

	if math.Abs(dy) > remainingRange {
		dy = math.Copysign(remainingRange, dy)
	}

	newX := c.location.X() + int(dx)
	newY := c.location.Y() + int(dy)

	newLocation, err := kernel.NewLocation(newX, newY)
	if err != nil {
		return err
	}
	c.location = newLocation
	return nil
}

func (c *Courier) findStoragePlaceByOrderID(orderID uuid.UUID) (*StoragePlace, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsInvalidError("orderID")
	}

	for _, place := range c.storagePlaces {
		if place.orderID != nil && *place.orderID == orderID {
			return place, nil
		}
	}

	return nil, ErrOrderNotFound
}
