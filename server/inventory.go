package server

import (
	inventoryGrpcHandler "microService/modules/inventory/inventoryHandler/grpc"
	inventoryHttpHandler "microService/modules/inventory/inventoryHandler/http"

	inventoryQueueHandler "microService/modules/inventory/inventoryHandler/inventoryQueue"

	"microService/modules/inventory/inventoryRepository"
	"microService/modules/inventory/inventoryUsecase"
)

func (s *server) inventpryService() {
	repo := inventoryRepository.NewInventoryRepository(s.mongo)
	usecase := inventoryUsecase.NewInventoryUsecase(repo)
	inventpryHttpHandler := inventoryHttpHandler.NewInventoryHttpHandler(s.cfg, usecase)
	inventpryGrpcHandler := inventoryGrpcHandler.NewInventoryGrpcHandler(usecase)
	inventpryQueueHandler := inventoryQueueHandler.NewInventoryQueueHandler(s.cfg, usecase)

	_ = inventpryHttpHandler
	_ = inventpryGrpcHandler
	_ = inventpryQueueHandler

	inventory := s.app.Group("/inventpry_v1")

	inventory.GET("/health", s.healthCheckService)

}
