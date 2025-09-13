package kernel

import (
	"delivery/internal/pkg/errs"
	"fmt"
	"math/rand/v2"
)

type Location struct {
	x       int
	y       int
	isValid bool
}

const (
	minXCoordinate = 1
	maxXCoordinate = 10
	minYCoordinate = 1
	maxYCoordinate = 10
)

func NewLocation(x, y int) (Location, error) {
	if x > maxXCoordinate || x < minXCoordinate {
		return Location{}, errs.NewValueIsInvalidError("x")
	}

	if y > maxYCoordinate || y < minYCoordinate {
		return Location{}, errs.NewValueIsInvalidError("y")
	}

	return Location{
		x:       x,
		y:       y,
		isValid: true,
	}, nil
}

func (l Location) X() int {
	return l.x
}

func (l Location) Y() int {
	return l.y
}

func (l Location) Equals(other Location) bool {
	return l == other
}

func (l Location) DistanceTo(other Location) int {
	return abs(l.x-other.x) + abs(l.y-other.y)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func RandomLocation() Location {
	loc, err := NewLocation(
		randomCoordinate(minXCoordinate, maxXCoordinate),
		randomCoordinate(minYCoordinate, maxYCoordinate),
	)

	if err != nil {
		panic(fmt.Errorf("incorrect random location: %w", err))
	}

	return loc
}

func MinLocation() Location {
	loc, err := NewLocation(minXCoordinate, minYCoordinate)
	if err != nil {
		panic(fmt.Errorf("incorrect min location: %w", err))
	}

	return loc
}

func MaxLocation() Location {
	loc, err := NewLocation(maxXCoordinate, maxYCoordinate)
	if err != nil {
		panic(fmt.Errorf("incorrect max location: %w", err))
	}

	return loc
}

func randomCoordinate(min, max int) int {
	return rand.IntN(max-min+1) + min
}

func (l Location) String() string {
	return fmt.Sprintf("(%d,%d)", l.x, l.y)
}

func (l Location) IsValid() bool {
	return l.isValid
}
