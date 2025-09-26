package cmd

import (
	"delivery/internal/adapters/in/kafka"
	"delivery/internal/adapters/out/grpc/geo"
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/jobs"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CompositionRoot struct {
	configs   Config
	gormDb    *gorm.DB
	geoClient ports.GeoClient
	onceGeo   sync.Once
	closers   []Closer
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
func (cr *CompositionRoot) NewCreateOrderCommandHandler() commands.CreateOrderCommandHandler {
	commandHandler, err := commands.NewCreateOrderCommandHandler(cr.NewUnitOfWorkFactory(), cr.NewGeoClient())
	if err != nil {
		log.Fatalf("cannot create CreateOrderCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewCreateCourierCommandHandler() commands.CreateCourierCommandHandler {
	commandHandler, err := commands.NewCreateCourierCommandHandler(cr.NewUnitOfWorkFactory())
	if err != nil {
		log.Fatalf("cannot create CreateCourierCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewAssignOrdersCommandHandler() commands.AssignOrdersCommandHandler {
	commandHandler, err := commands.NewAssignOrdersCommandHandler(
		cr.NewUnitOfWorkFactory(), cr.NewOrderDispatcher())
	if err != nil {
		log.Fatalf("cannot create AssignOrdersCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewMoveCouriersCommandHandler() commands.MoveCouriersCommandHandler {
	commandHandler, err := commands.NewMoveCouriersCommandHandler(
		cr.NewUnitOfWorkFactory())
	if err != nil {
		log.Fatalf("cannot create MoveCouriersCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewGetAllCouriersQueryHandler() queries.GetAllCouriersQueryHandler {
	queryHandler, err := queries.NewGetAllCouriersQueryHandler(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create GetAllCouriersQueryHandler: %v", err)
	}
	return queryHandler
}

func (cr *CompositionRoot) NewGetNotCompletedOrdersQueryHandler() queries.GetNotCompletedOrdersQueryHandler {
	queryHandler, err := queries.NewGetNotCompletedOrdersQueryHandler(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create GetNotCompletedOrdersQueryHandler: %v", err)
	}
	return queryHandler
}

func (cr *CompositionRoot) NewAssignOrdersJob() cron.Job {
	job, err := jobs.NewAssignOrdersJob(cr.NewAssignOrdersCommandHandler())
	if err != nil {
		log.Fatalf("cannot create AssignOrdersJob: %v", err)
	}
	return job
}

func (cr *CompositionRoot) NewMoveCouriersJob() cron.Job {
	job, err := jobs.NewMoveCouriersJob(cr.NewMoveCouriersCommandHandler())
	if err != nil {
		log.Fatalf("cannot create MoveCouriersJob: %v", err)
	}
	return job
}

func (cr *CompositionRoot) NewGeoClient() ports.GeoClient {
	cr.onceGeo.Do(func() {
		client, err := geo.NewClient(cr.configs.GeoServiceGrpcHost)
		if err != nil {
			log.Fatalf("cannot create GeoClient: %v", err)
		}

		cr.RegisterCloser(client)
		cr.geoClient = client
	})
	return cr.geoClient
}

func (cr *CompositionRoot) NewBasketConfirmedConsumer() kafka.BasketConfirmedConsumer {
	consumer, err := kafka.NewBasketConfirmedConsumer(
		[]string{cr.configs.KafkaHost},
		cr.configs.KafkaConsumerGroup,
		cr.configs.KafkaBasketConfirmedTopic,
		cr.NewCreateOrderCommandHandler(),
	)

	if err != nil {
		log.Fatalf("cannot create BasketConfirmedConsumer: %v", err)
	}

	cr.RegisterCloser(consumer)

	return consumer
}
