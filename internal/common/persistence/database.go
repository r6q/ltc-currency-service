package persistence

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/robertst0/ltc-currency-service/internal/common/config"
	"github.com/robertst0/ltc-currency-service/internal/component/currency"
)

func NewDatabase(properties config.DatabaseProperties) *gorm.DB {
	db, err := gorm.Open(mysql.Open(properties.ConnectionURL))
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	err = db.AutoMigrate(currency.RateEntity{})
	if err != nil {
		log.Fatal("failed to migrate", err)
	}

	return db
}
