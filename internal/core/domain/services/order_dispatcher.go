package services

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"errors"
	"math"
)

var (
	ErrOrderIsAlreadyAssigned = errors.New("order is already assigned")
	ErrNoSuitableCourier      = errors.New("no suitable courier")
)

type OrderDispatcher interface {
	Dispatch(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error)
}

var _ OrderDispatcher = &orderDispatcher{}

type orderDispatcher struct{}

func NewOrderDispatcher() OrderDispatcher {
	return &orderDispatcher{}
}

func (od orderDispatcher) Dispatch(o *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if o == nil {
		return nil, errs.NewValueIsRequiredError("order")
	}

	if len(couriers) == 0 {
		return nil, errs.NewValueIsRequiredError("couriers")
	}

	if o.Status() != order.StatusCreated {
		return nil, ErrOrderIsAlreadyAssigned
	}

	bestCourier, err := od.selectBestCourier(o, couriers)
	if err != nil {
		return nil, err
	}

	if bestCourier == nil {
		return nil, ErrNoSuitableCourier
	}

	err = bestCourier.TakeOrder(o)
	if err != nil {
		return nil, err
	}

	err = o.Assign(bestCourier.Id())
	if err != nil {
		return nil, err
	}

	return bestCourier, nil
}

func (od orderDispatcher) selectBestCourier(o *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	var (
		bestCourier *courier.Courier
		minTime     = math.MaxFloat64
	)

	for _, c := range couriers {
		canTake, err := c.CanTakeOrder(o)
		if err != nil {
			return nil, err
		}

		if !canTake {
			continue
		}

		t, err := c.CalculateTimeToLocation(o.Location())
		if err != nil {
			return nil, err
		}

		if t < minTime {
			minTime = t

			bestCourier = c
		}
	}

	return bestCourier, nil
}
