package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Route(group *echo.Group) {
	group.GET("", h.findAll())
}

func (h Handler) findAll() func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		if ctx.QueryParam("latest") != "" {
			return h.findLatest(ctx)
		}

		if ctx.QueryParam("currency") != "" {
			return h.findHistorical(ctx)
		}

		return ctx.JSON(http.StatusOK, DataResponse{h.service.FindAll()})
	}
}

func (h *Handler) findLatest(ctx echo.Context) error {
	param := ctx.QueryParam("latest")

	latest, err := strconv.ParseBool(param)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{fmt.Sprintf("Param value '%s' is invalid", param)})
	}

	if latest {
		return ctx.JSON(http.StatusOK, DataResponse{h.service.FindLatest()})
	}

	return ctx.JSON(http.StatusOK, DataResponse{h.service.FindAll()})
}

func (h *Handler) findHistorical(ctx echo.Context) error {
	currency := ctx.QueryParam("currency")

	return ctx.JSON(http.StatusOK, DataResponse{h.service.FindHistorical(currency)})
}
