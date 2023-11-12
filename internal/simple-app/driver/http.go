// Package driver contains driver layer
package driver

import (
	"errors"
	"net/http"

	"github.com/adlandh/acorn-simple-app/internal/simple-app/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

//go:generate oapi-codegen -old-config-style -generate types,server -o "openapi_gen.go" -package "driver" "../../../api/simple-app.yaml"
type HTTPServer struct {
	app domain.ApplicationInterface
}

var _ ServerInterface = (*HTTPServer)(nil)

func NewHTTPServer(app domain.ApplicationInterface) *HTTPServer {
	return &HTTPServer{
		app: app,
	}
}

func (h HTTPServer) HealthCheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Ok")
}

func (h HTTPServer) CreateUser(ctx echo.Context) error {
	var userRequest UserRequest

	err := ctx.Bind(&userRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := h.app.CreateUser(ctx.Request().Context(), userRequest.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, User{
		Id:   id,
		Name: userRequest.Name,
	})
}

func (h HTTPServer) DeleteUser(ctx echo.Context, id uuid.UUID) error {
	err := h.app.DeleteUser(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(domain.ErrorNotFound, err) {
			return ctx.NoContent(http.StatusNotFound)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (h HTTPServer) GetUser(ctx echo.Context, id uuid.UUID) error {
	name, err := h.app.GetUser(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(domain.ErrorNotFound, err) {
			return ctx.NoContent(http.StatusNotFound)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, User{
		Id:   id,
		Name: name,
	})
}

func (h HTTPServer) UpdateUser(ctx echo.Context, id uuid.UUID) error {
	var userRequest UserRequest

	err := ctx.Bind(&userRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.app.UpdateUser(ctx.Request().Context(), id, userRequest.Name)

	if err != nil {
		if errors.Is(domain.ErrorNotFound, err) {
			return ctx.NoContent(http.StatusNotFound)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, User{
		Id:   id,
		Name: userRequest.Name,
	})
}
