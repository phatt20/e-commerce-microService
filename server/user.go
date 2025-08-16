package server

import (
	"log"
	userGrpcHandler "microService/modules/user/userHandler/grpc"
	userHttpHandler "microService/modules/user/userHandler/http"
	userQueueHandler "microService/modules/user/userHandler/userQueue"
	"microService/modules/user/userPb"
	grpccon "microService/pkg/grpcCon"

	"microService/modules/user/userRepository"
	"microService/modules/user/userUsecase"
)

func (s *server) userService() {
	repo := userRepository.NewUserRepository(s.mongo)
	usecase := userUsecase.NewUserUsecase(repo)
	userHttpHandler := userHttpHandler.NewUserHttpHandler(s.cfg, usecase)
	userGrpcHandler := userGrpcHandler.NewUserGrpcHandler(usecase)
	userQueueHandler := userQueueHandler.NewUserQueueHandler(s.cfg, usecase)

	go userQueueHandler.DockedUserMoney()
	go userQueueHandler.AddUserMoney()
	go userQueueHandler.RollbackUserTransaction()

	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.UserUrl)

		userPb.RegisterUserGrpcServiceServer(grpcServer, userGrpcHandler)

		log.Printf("User gRPC server listening on %s", s.cfg.Grpc.UserUrl)
		grpcServer.Serve(lis)
	}()

	user := s.app.Group("/user_v1")

	user.POST("/user/register", userHttpHandler.CreateUser)
	user.POST("/user/add-money", userHttpHandler.AddUserMoney, s.middleware.JwtAuthorization)
	user.GET("/user/:user_id", userHttpHandler.FindOneUserProfile)
	user.GET("/user/saving-account/my-account", userHttpHandler.GetUserSavingAccount, s.middleware.JwtAuthorization)
}
