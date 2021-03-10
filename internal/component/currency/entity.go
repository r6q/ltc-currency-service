package currency

import "time"

const RateTable = "rate"

type RateEntity struct {
	ID           int64  `gorm:"primaryKey"`
	CurrencyCode string `gorm:"index:currency_by_date,unique"`
	Amount       string
	Date         time.Time `gorm:"index:currency_by_date,unique"`
}

func (RateEntity) TableName() string {
	return RateTable
}
