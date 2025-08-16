package server

import (
	itemGrpcHandler "microService/modules/item/itemHandler/grpc"
	itemHttpHandler "microService/modules/item/itemHandler/http"
	"microService/modules/item/itemRepository"
	"microService/modules/item/itemUsecase"
)

func (s *server) itemService() {
	repo := itemRepository.NewitemRepository(s.mongo)
	usecase := itemUsecase.NewitemUsecase(repo)
	itemHttpHandler := itemHttpHandler.NewitemHttpHandler(s.cfg, usecase)
	itemGrpcHandler := itemGrpcHandler.NewitemGrpcHandler(usecase)

	_ = itemHttpHandler
	_ = itemGrpcHandler

	item := s.app.Group("/item_v1")

	item.GET("/health", s.healthCheckService)

}
