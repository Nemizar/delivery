package tests

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"

	"github.com/google/uuid"
)

func CreateLocation(x, y int) kernel.Location {
	l, err := kernel.NewLocation(x, y)
	if err != nil {
		panic(err)
	}

	return l
}

func CreateCourier(name string, speed int, location kernel.Location) *courier.Courier {
	c, err := courier.NewCourier(name, speed, location)
	if err != nil {
		panic(err)
	}

	return c
}

func CreateOrder(uuid uuid.UUID, location kernel.Location, volume int) *order.Order {
	o, err := order.NewOrder(uuid, location, volume)
	if err != nil {
		panic(err)
	}

	return o
}
