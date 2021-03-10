package currency

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Persist(rate *RateEntity) {
	r.db.Clauses(clause.OnConflict{
		OnConstraint: "currency_by_date",
		DoUpdates:    clause.AssignmentColumns([]string{"amount"}),
	}).Create(rate)
}

func (r *Repository) FindAll() []*RateEntity {
	var rates []*RateEntity

	r.db.Find(&rates)

	return rates
}

func (r *Repository) FindLatest() []*RateEntity {
	var rates []*RateEntity

	r.db.Where("date = (?)", r.db.Table(RateTable).Select("MAX(date)")).Find(&rates)

	return rates
}

func (r *Repository) FindByCurrency(currency string) []*RateEntity {
	var rates []*RateEntity

	r.db.Find(&rates, "currency_code = ?", currency)

	return rates
}
