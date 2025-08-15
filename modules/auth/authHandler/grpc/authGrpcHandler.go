package authHandler

import (
	"context"
	"microService/modules/auth/authPb"
	authusecase "microService/modules/auth/authUsecase"
)

type (
	authGrpcHandler struct {
		authPb.UnimplementedAuthGrpcServiceServer
		authUsecase authusecase.AuthUsecaseService
	}
)

func AuthGrpcHandler(authUsecase authusecase.AuthUsecaseService) *authGrpcHandler {
	return &authGrpcHandler{authUsecase: authUsecase}
}

func (g *authGrpcHandler) AccessTokenSearch(ctx context.Context, req *authPb.AccessTokenSearchReq) (*authPb.AccessTokenSearchRes, error) {
	return g.authUsecase.AccessTokenSearch(ctx, req.AccessToken)
}

func (g *authGrpcHandler) RolesCount(ctx context.Context, req *authPb.RolesCountReq) (*authPb.RolesCountRes, error) {
	return g.authUsecase.RolesCount(ctx)
}
