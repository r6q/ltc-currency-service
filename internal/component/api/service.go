package api

import (
	"github.com/robertst0/ltc-currency-service/internal/component/currency"
)

type Service struct {
	repository *currency.Repository
}

func NewService(repository *currency.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) FindAll() []*Rate {
	return s.remap(s.repository.FindAll())
}

func (s *Service) FindLatest() []*Rate {
	return s.remap(s.repository.FindLatest())
}

func (s *Service) FindHistorical(currency string) []*Rate {
	return s.remap(s.repository.FindByCurrency(currency))
}

func (s *Service) remap(rateEntities []*currency.RateEntity) []*Rate {
	rates := make([]*Rate, len(rateEntities))
	for i, entity := range rateEntities {
		rates[i] = &Rate{
			CurrencyCode: entity.CurrencyCode,
			Amount:       entity.Amount,
			Date:         entity.Date,
		}
	}

	return rates
}
