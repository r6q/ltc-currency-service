package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/robertst0/ltc-currency-service/internal/common/config"
	"github.com/robertst0/ltc-currency-service/internal/common/persistence"
	"github.com/robertst0/ltc-currency-service/internal/component/api"
	"github.com/robertst0/ltc-currency-service/internal/component/currency"
)

func main() {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	properties := config.LoadProperties()

	repository := currency.NewRepository(persistence.NewDatabase(properties.Database))

	handler := api.NewHandler(api.NewService(repository))
	handler.Route(app.Group("/api/v1/rates"))

	app.Logger.Fatal(app.Start(":8080"))
}
