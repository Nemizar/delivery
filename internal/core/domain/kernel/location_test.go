package kernel_test

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/pkg/errs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLocation(t *testing.T) {
	type args struct {
		x int
		y int
	}

	tests := []struct {
		name    string
		args    args
		wantX   int
		wantY   int
		wantErr error
	}{
		{
			name: "Valid coordinate values",
			args: args{
				x: 2,
				y: 5,
			},
			wantX: 2,
			wantY: 5,
		},
		{
			name: "Min coordinate values",
			args: args{
				x: 1,
				y: 1,
			},
			wantX: 1,
			wantY: 1,
		},
		{
			name: "Max coordinate values",
			args: args{
				x: 10,
				y: 10,
			},
			wantX: 10,
			wantY: 10,
		},
		{
			name: "Negative values for coordinates",
			args: args{
				x: -1,
				y: -1,
			},
			wantErr: errs.NewValueIsInvalidError("x"),
		},
		{
			name: "Negative value for X only",
			args: args{
				x: -1,
				y: 5,
			},
			wantErr: errs.NewValueIsInvalidError("x"),
		},
		{
			name: "Negative value for Y only",
			args: args{
				x: 5,
				y: -1,
			},
			wantErr: errs.NewValueIsInvalidError("y"),
		},
		{
			name: "Zero coordinate values",
			args: args{
				x: 0,
				y: 0,
			},
			wantErr: errs.NewValueIsInvalidError("x"),
		},
		{
			name: "Coordinate values exceeding limits",
			args: args{
				x: 100,
				y: 100,
			},
			wantErr: errs.NewValueIsInvalidError("x"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := kernel.NewLocation(tt.args.x, tt.args.y)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Equal(t, kernel.Location{}, l)

				return
			}

			require.Nil(t, err)

			assert.Equal(t, tt.wantX, l.X())
			assert.Equal(t, tt.wantY, l.Y())
		})
	}
}

func TestLocation_Equals(t *testing.T) {
	tests := []struct {
		name   string
		first  kernel.Location
		second kernel.Location
		want   bool
	}{
		{
			name:   "Equals",
			first:  createValidLocation(1, 1),
			second: createValidLocation(1, 1),
			want:   true,
		},
		{
			name:   "Not equals 1",
			first:  createValidLocation(1, 2),
			second: createValidLocation(1, 1),
		},
		{
			name:   "Not equals 2",
			first:  createValidLocation(1, 1),
			second: createValidLocation(1, 2),
		},
		{
			name:   "Not equals 3",
			first:  createValidLocation(2, 1),
			second: createValidLocation(1, 2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.first.Equals(tt.second))
			assert.Equal(t, tt.want, tt.second.Equals(tt.first))
			assert.True(t, tt.first.Equals(tt.first))
			assert.True(t, tt.second.Equals(tt.second))
		})
	}
}

func TestLocation_String(t *testing.T) {
	l := createValidLocation(3, 7)

	assert.Equal(t, "(3,7)", l.String())
}

func TestLocation_DistanceTo(t *testing.T) {
	tests := []struct {
		name            string
		courierLocation kernel.Location
		orderLocation   kernel.Location
		want            int
	}{
		{
			name:            "5 steps",
			courierLocation: createValidLocation(4, 9),
			orderLocation:   createValidLocation(2, 6),
			want:            5,
		},
		{
			name:            "Min steps",
			courierLocation: createValidLocation(4, 4),
			orderLocation:   createValidLocation(4, 4),
			want:            0,
		},
		{
			name:            "Max steps",
			courierLocation: createValidLocation(1, 1),
			orderLocation:   createValidLocation(10, 10),
			want:            18,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.courierLocation.DistanceTo(tt.orderLocation))
			assert.Equal(t, tt.want, tt.orderLocation.DistanceTo(tt.courierLocation))
		})
	}
}

func TestRandomLocation_Range(t *testing.T) {
	for i := 0; i < 100; i++ {
		loc := kernel.RandomLocation()
		assert.GreaterOrEqual(t, loc.X(), 1)
		assert.LessOrEqual(t, loc.X(), 10)
		assert.GreaterOrEqual(t, loc.Y(), 1)
		assert.LessOrEqual(t, loc.Y(), 10)
	}
}

func TestRandomLocation_Variety(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		loc := kernel.RandomLocation()
		seen[loc.String()] = true
	}

	assert.Greater(t, len(seen), 10)
}

func createValidLocation(x, y int) kernel.Location {
	l, err := kernel.NewLocation(x, y)
	if err != nil {
		panic(err)
	}

	return l
}
