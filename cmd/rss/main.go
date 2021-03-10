package main

import (
	"github.com/robertst0/ltc-currency-service/internal/common/config"
	"github.com/robertst0/ltc-currency-service/internal/common/persistence"
	"github.com/robertst0/ltc-currency-service/internal/component/currency"
	"github.com/robertst0/ltc-currency-service/internal/component/rss"
)

func main() {
	properties := config.LoadProperties()

	service := rss.NewService(
		currency.NewRepository(persistence.NewDatabase(properties.Database)),
		properties.Rss,
	)

	service.RetrieveLatestData()
}
