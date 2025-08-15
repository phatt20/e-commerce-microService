package userHandler

import (
	"context"
	userpb "microService/modules/user/userPb"
	"microService/modules/user/userUsecase"
)

type (
	userGrpcHandler struct {
		userpb.UnimplementedUserGrpcServiceServer
		userUsecase userUsecase.UserUsecase
	}
)

func NewUserGrpcHandler(userUsecase userUsecase.UserUsecase) *userGrpcHandler {
	return &userGrpcHandler{userUsecase: userUsecase}
}
func (g *userGrpcHandler) CredentialSearch(ctx context.Context, req *userpb.CredentialSearchReq) (*userpb.UserProfile, error) {
	return g.userUsecase.FindOneUserCredential(ctx, req.Password, req.Email)
}

func (g *userGrpcHandler) FindOneUserProfileToRefresh(ctx context.Context, req *userpb.FindOneUserProfileToRefreshReq) (*userpb.UserProfile, error) {
	return g.userUsecase.FindOneUserProfileToRefresh(ctx, req.UserId)
}

func (g *userGrpcHandler) GetUserSavingAccount(ctx context.Context, req *userpb.GetUserSavingAccountReq) (*userpb.GetUserSavingAccountRes, error) {
	return nil, nil
}
