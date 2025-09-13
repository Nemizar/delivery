package cmd

import (
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"log"

	"gorm.io/gorm"
)

type CompositionRoot struct {
	configs Config
	gormDb  *gorm.DB

	closers []Closer
}

func NewCompositionRoot(configs Config, gormDb *gorm.DB) *CompositionRoot {
	return &CompositionRoot{
		configs: configs,
		gormDb:  gormDb,
	}
}

func (cr *CompositionRoot) NewOrderDispatcher() services.OrderDispatcher {
	orderDispatcher := services.NewOrderDispatcher()
	return orderDispatcher
}

func (cr *CompositionRoot) NewUnitOfWork() ports.UnitOfWork {
	unitOfWork, err := postgres.NewUnitOfWork(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}
	return unitOfWork
}

func (cr *CompositionRoot) NewUnitOfWorkFactory() ports.UnitOfWorkFactory {
	unitOfWorkFactory, err := postgres.NewUnitOfWorkFactory(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create UnitOfWorkFactory: %v", err)
	}
	return unitOfWorkFactory
}
