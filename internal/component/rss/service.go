package rss

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/robertst0/ltc-currency-service/internal/common/config"
	"github.com/robertst0/ltc-currency-service/internal/component/currency"
)

type Service struct {
	repository *currency.Repository
	properties config.RssProperties
	client     http.Client
}

func NewService(repository *currency.Repository, properties config.RssProperties) *Service {
	return &Service{
		repository: repository,
		properties: properties,
		client:     http.Client{Timeout: properties.RequestTimeout},
	}
}

func (s *Service) RetrieveLatestData() {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, s.properties.URL, nil)
	if err != nil {
		log.Println(err)

		return
	}

	response, err := s.client.Do(req)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Println(err)

		return
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)

		return
	}

	var data Response

	err = xml.Unmarshal(bytes, &data)
	if err != nil {
		log.Println(err)

		return
	}

	for item := range data.Channel.Items {
		item := &data.Channel.Items[item]
		for rateIndex := range item.Currencies.Rates {
			rate := &item.Currencies.Rates[rateIndex]

			s.repository.Persist(&currency.RateEntity{
				CurrencyCode: rate.Code,
				Amount:       rate.Amount,
				Date:         item.Date.Time,
			})
		}
	}

	log.Println("Successfully updated exchange rates")

	response.Body.Close()
}
