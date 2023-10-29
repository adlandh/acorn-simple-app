package main

import (
	"github.com/adlandh/acorn-simple-app/internal/simple-app/application"
	"github.com/adlandh/acorn-simple-app/internal/simple-app/config"
	"github.com/adlandh/acorn-simple-app/internal/simple-app/domain"
	"github.com/adlandh/acorn-simple-app/internal/simple-app/driven"
	"github.com/adlandh/acorn-simple-app/internal/simple-app/driver"

	"context"
	"errors"
	"fmt"
	"net/http"

	echoZapMiddleware "github.com/adlandh/echo-zap-middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(createService()).Run()
}

func createService() fx.Option {
	return fx.Options(
		fx.WithLogger(
			func(log *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: log}
			},
		),
		fx.Provide(
			config.NewConfig,
			fx.Annotate(
				zap.NewDevelopment,
			),
			fx.Annotate(
				driven.NewRedisStorage,
				fx.As(new(domain.UserStorage)),
			),
			fx.Annotate(
				application.NewApplication,
				fx.As(new(domain.ApplicationInterface)),
			),
			fx.Annotate(
				driver.NewHTTPServer,
				fx.As(new(driver.ServerInterface)),
			),
		),
		fx.Invoke(
			newEcho,
		),
	)
}

func newEcho(lc fx.Lifecycle, server driver.ServerInterface, cfg *config.Config, log *zap.Logger) *echo.Echo {
	e := echo.New()
	e.Use(echoZapMiddleware.Middleware(log))
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("1M"))
	e.Use(middleware.RequestID())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			driver.RegisterHandlers(e, server)
			go func() {
				err = e.Start(":" + cfg.Port)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("error starting echo server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := e.Shutdown(ctx)
			if err != nil {
				return fmt.Errorf("error shutting down echo server: %w", err)
			}
			return nil
		},
	})

	return e
}
