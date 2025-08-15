package authusecase

import (
	"context"
	"errors"
	"log"
	"microService/config"
	"microService/modules/auth"
	"microService/modules/auth/authPb"
	authrepository "microService/modules/auth/authRepository"
	"microService/modules/user/userPb"
	"microService/pkg/jwtauth"
	"microService/pkg/utils"
	"strings"
	"time"
)

type (
	AuthUsecaseService interface {
		Login(pctx context.Context, cfg *config.Config, req *auth.UserLoginReq) (*auth.ProfileIntercepter, error)
		RefreshToken(pctx context.Context, cfg *config.Config, req *auth.RefreshTokenReq) (*auth.ProfileIntercepter, error)
		Logout(pctx context.Context, credentialId string) (int64, error)
		AccessTokenSearch(pctx context.Context, accessToken string) (*authPb.AccessTokenSearchRes, error)
		RolesCount(pctx context.Context) (*authPb.RolesCountRes, error)
	}

	authUsecase struct {
		authRepository authrepository.AuthRepositoryService
	}
)

func NewAuthUsecase(authRepository authrepository.AuthRepositoryService) AuthUsecaseService {
	return &authUsecase{authRepository}
}

func (u *authUsecase) Login(pctx context.Context, cfg *config.Config, req *auth.UserLoginReq) (*auth.ProfileIntercepter, error) {
	profile, err := u.authRepository.CredentialSearch(pctx, cfg.Grpc.UserUrl, &userPb.CredentialSearchReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	profile.Id = "user:" + profile.Id

	accessToken := u.authRepository.AccessToken(cfg, &jwtauth.Claims{
		UserId:   profile.Id,
		RoleCode: int(profile.RoleCode),
	})

	refreshToken := u.authRepository.RefreshToken(cfg, &jwtauth.Claims{
		UserId:   profile.Id,
		RoleCode: int(profile.RoleCode),
	})

	credentialId, err := u.authRepository.InsertOneUserCredential(pctx, &auth.Credential{
		UserId:       profile.Id,
		RoleCode:     int(profile.RoleCode),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	credential, err := u.authRepository.FindOneUserCredential(pctx, credentialId.Hex())
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &auth.ProfileIntercepter{
		UserProfile: &userPb.UserProfile{
			Id:        profile.Id,
			Email:     profile.Email,
			Username:  profile.Username,
			CreatedAt: utils.ConvertStringTimeToTime(profile.CreatedAt).In(loc).Format(time.RFC3339),
			UpdatedAt: utils.ConvertStringTimeToTime(profile.UpdatedAt).In(loc).Format(time.RFC3339),
		},
		Credential: &auth.CredentialRes{
			Id:           credential.Id.Hex(),
			UserId:       credential.UserId,
			RoleCode:     credential.RoleCode,
			AccessToken:  credential.AccessToken,
			RefreshToken: credential.RefreshToken,
			CreatedAt:    credential.CreatedAt.In(loc),
			UpdatedAt:    credential.UpdatedAt.In(loc),
		},
	}, nil
}

func (u *authUsecase) RefreshToken(pctx context.Context, cfg *config.Config, req *auth.RefreshTokenReq) (*auth.ProfileIntercepter, error) {
	claims, err := jwtauth.ParseToken(cfg.Jwt.RefreshSecretKey, req.RefreshToken)
	if err != nil {
		log.Printf("Error: RefreshToken: %s", err.Error())
		return nil, errors.New(err.Error())
	}

	profile, err := u.authRepository.FindOneUserProfileToRefresh(pctx, cfg.Grpc.UserUrl, &userPb.FindOneUserProfileToRefreshReq{
		UserId: strings.TrimPrefix(claims.UserId, "user:"),
	})
	if err != nil {
		return nil, err
	}

	accessToken := jwtauth.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuration, &jwtauth.Claims{
		UserId:   profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()

	refreshToken := jwtauth.ReloadToken(cfg.Jwt.RefreshSecretKey, claims.ExpiresAt.Unix(), &jwtauth.Claims{
		UserId:   profile.Id,
		RoleCode: int(profile.RoleCode),
	})

	if err := u.authRepository.UpdateOneUserCredential(pctx, req.CredentialId, &auth.UpdateRefreshTokenReq{
		UserId:       profile.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UpdatedAt:    utils.LocalTime(),
	}); err != nil {
		return nil, err
	}

	credential, err := u.authRepository.FindOneUserCredential(pctx, req.CredentialId)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &auth.ProfileIntercepter{
		UserProfile: &userPb.UserProfile{
			Id:        "user:" + profile.Id,
			Email:     profile.Email,
			Username:  profile.Username,
			CreatedAt: utils.ConvertStringTimeToTime(profile.CreatedAt).Format(time.RFC3339),
			UpdatedAt: utils.ConvertStringTimeToTime(profile.UpdatedAt).Format(time.RFC3339),
		},
		Credential: &auth.CredentialRes{
			Id:           credential.Id.Hex(),
			UserId:       credential.UserId,
			RoleCode:     credential.RoleCode,
			AccessToken:  credential.AccessToken,
			RefreshToken: credential.RefreshToken,
			CreatedAt:    credential.CreatedAt.In(loc),
			UpdatedAt:    credential.UpdatedAt.In(loc),
		},
	}, nil
}

func (u *authUsecase) Logout(pctx context.Context, credentialId string) (int64, error) {
	return u.authRepository.DeleteOneUserCredential(pctx, credentialId)
}

func (u *authUsecase) AccessTokenSearch(pctx context.Context, accessToken string) (*authPb.AccessTokenSearchRes, error) {
	credential, err := u.authRepository.FindOneAccessToken(pctx, accessToken)
	if err != nil {
		return &authPb.AccessTokenSearchRes{
			IsValid: false,
		}, err
	}

	if credential == nil {
		return &authPb.AccessTokenSearchRes{
			IsValid: false,
		}, errors.New("error: access token is invalid")
	}

	return &authPb.AccessTokenSearchRes{
		IsValid: true,
	}, nil
}

func (u *authUsecase) RolesCount(pctx context.Context) (*authPb.RolesCountRes, error) {
	result, err := u.authRepository.RolesCount(pctx)
	if err != nil {
		return nil, err
	}

	return &authPb.RolesCountRes{
		Count: result,
	}, nil
}
