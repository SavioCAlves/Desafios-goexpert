//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/SavioCAlves/Desafios-goexpert/tree/main/Clean_Architecture/internal/entity"
	"github.com/SavioCAlves/Desafios-goexpert/tree/main/Clean_Architecture/internal/event"
	"github.com/SavioCAlves/Desafios-goexpert/tree/main/Clean_Architecture/internal/infra/database"
	"github.com/SavioCAlves/Desafios-goexpert/tree/main/Clean_Architecture/internal/infra/web"
	"github.com/SavioCAlves/Desafios-goexpert/tree/main/Clean_Architecture/internal/usecase"
	"github.com/SavioCAlves/Desafios-goexpert/tree/main/Clean_Architecture/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
