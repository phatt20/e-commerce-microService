package middlewarehandler

import (
	"microService/config"
	middlewareusecase "microService/modules/middleware/middlewareUsecase"
	"microService/pkg/response"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	MiddlewareHandler interface {
		JwtAuthorization(next echo.HandlerFunc) echo.HandlerFunc
		RbacAuthorization(next echo.HandlerFunc, expected []int) echo.HandlerFunc
		UserIdParamValidation(next echo.HandlerFunc) echo.HandlerFunc
	}

	middlewareHandler struct {
		cfg               *config.Config
		middlewareUsecase middlewareusecase.MiddlewareUsecase
	}
)

func NewMiddlewareHandler(cfg *config.Config, middlewareusecase middlewareusecase.MiddlewareUsecase) MiddlewareHandler {
	return &middlewareHandler{cfg, middlewareusecase}
}

func (h *middlewareHandler) JwtAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")

		newCtx, err := h.middlewareUsecase.JwtAuthorization(c, h.cfg, accessToken)
		if err != nil {
			return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
		}

		return next(newCtx)
	}
}

func (h *middlewareHandler) RbacAuthorization(next echo.HandlerFunc, expected []int) echo.HandlerFunc {
	return func(c echo.Context) error {
		newCtx, err := h.middlewareUsecase.RbacAuthorization(c, h.cfg, expected)
		if err != nil {
			return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
		}

		return next(newCtx)
	}
}

func (h *middlewareHandler) UserIdParamValidation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		newCtx, err := h.middlewareUsecase.UserIdParamValidation(c)
		if err != nil {
			return response.ErrResponse(c, http.StatusUnauthorized, err.Error())
		}

		return next(newCtx)
	}
}
