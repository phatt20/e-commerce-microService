package server

import (
	"log"
	authGrpcHandler "microService/modules/auth/authHandler/grpc"
	authHttpHandler "microService/modules/auth/authHandler/http"
	"microService/modules/auth/authPb"
	grpccon "microService/pkg/grpcCon"

	authrepository "microService/modules/auth/authRepository"
	authusecase "microService/modules/auth/authUsecase"
)

func (s *server) authService() {
	repo := authrepository.NewAuthRepository(s.db)
	usecase := authusecase.NewAuthUsecase(repo)
	httpHandler := authHttpHandler.NewAuthHttpHandler(s.cfg, usecase)
	authGrpcHandler := authGrpcHandler.AuthGrpcHandler(usecase)

	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)

		authPb.RegisterAuthGrpcServiceServer(grpcServer, authGrpcHandler)

		log.Printf("Auth gRPC server listening on %s", s.cfg.Grpc.AuthUrl)
		grpcServer.Serve(lis)
	}()

	auth := s.app.Group("/auth_v1")

	auth.GET("", s.healthCheckService)
	auth.GET("/test/:user_id", s.healthCheckService)
	auth.POST("/auth/login", httpHandler.Login)
	auth.POST("/auth/refresh-token", httpHandler.RefreshToken)
	auth.POST("/auth/logout", httpHandler.Logout)
}
