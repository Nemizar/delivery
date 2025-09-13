package postgres

import (
	"context"
	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/adapters/out/postgres/orderrepo"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/testcnts"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (context.Context, *gorm.DB, error) {
	ctx := context.Background()
	postgresContainer, dsn, err := testcnts.StartPostgresContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	db, err := gorm.Open(postgresgorm.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&courierrepo.CourierDTO{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&courierrepo.StoragePlaceDTO{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&orderrepo.OrderDTO{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&courierrepo.StoragePlaceDTO{})
	assert.NoError(t, err)

	t.Cleanup(func() {
		err := postgresContainer.Terminate(ctx)
		assert.NoError(t, err)
	})

	return ctx, db, nil
}

func TestUnitOfWork_CourierRepositoryShouldCanAddCourier(t *testing.T) {
	ctx, db, err := setupTest(t)
	require.Nil(t, err)

	uow, err := NewUnitOfWork(db)
	require.Nil(t, err)

	location := kernel.MaxLocation()
	courierAggregate, err := courier.NewCourier("Велосипедист", 2, location)
	err = uow.CourierRepository().Add(ctx, courierAggregate)
	assert.NoError(t, err)

	var courierFromDb courierrepo.CourierDTO
	err = db.First(&courierFromDb, "id = ?", courierAggregate.Id()).Error
	assert.NoError(t, err)

	assert.Equal(t, courierAggregate.Id(), courierFromDb.ID)
	assert.Equal(t, courierAggregate.Speed(), courierFromDb.Speed)
}

func TestUnitOfWork_CourierRepositoryGetAllFree(t *testing.T) {
	ctx, db, err := setupTest(t)
	require.Nil(t, err)

	uow, err := NewUnitOfWork(db)
	require.Nil(t, err)

	courier1, err := courier.NewCourier("Велосипедист", 2, kernel.MaxLocation())
	assert.Nil(t, err)

	courier2, err := courier.NewCourier("Велосипедист", 2, kernel.MaxLocation())
	assert.Nil(t, err)

	err = uow.CourierRepository().Add(ctx, courier1)
	assert.Nil(t, err)

	err = uow.CourierRepository().Add(ctx, courier2)
	assert.Nil(t, err)

	got, err := uow.CourierRepository().GetAllFree(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(got))
}

func TestUnitOfWork_OrderRepositoryShouldCanAddOrder(t *testing.T) {
	ctx, db, err := setupTest(t)
	require.Nil(t, err)

	uow, err := NewUnitOfWork(db)
	require.Nil(t, err)

	location := kernel.MinLocation()
	orderAggregate, err := order.NewOrder(uuid.New(), location, 10)
	err = uow.OrderRepository().Add(ctx, orderAggregate)
	assert.NoError(t, err)

	var orderFromDb orderrepo.OrderDTO
	err = db.First(&orderFromDb, "id = ?", orderAggregate.ID()).Error
	assert.NoError(t, err)

	assert.Equal(t, orderAggregate.ID(), orderFromDb.ID)
	assert.Equal(t, orderAggregate.Location().X(), orderFromDb.Location.X)
	assert.Equal(t, orderAggregate.Location().Y(), orderFromDb.Location.Y)
	assert.Equal(t, orderAggregate.Volume(), orderFromDb.Volume)
	assert.Equal(t, orderAggregate.Status(), orderFromDb.Status)
}
