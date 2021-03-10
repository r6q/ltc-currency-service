package rss_test

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/robertst0/ltc-currency-service/internal/component/rss"
)

func TestDateTime_UnmarshalXML(t *testing.T) {
	xmlResponse := `<pubDate>Mon, 01 Mar 2021 02:00:00 +0200</pubDate>`

	var input rss.DateTime

	assertions := assert.New(t)
	assertions.NoError(xml.Unmarshal([]byte(xmlResponse), &input))
	assertions.Equal(time.Date(2021, time.March, 1, 0, 0, 0, 0, time.UTC), input.Time)
}

func TestCurrency_UnmarshalXML(t *testing.T) {
	xmlResponse := `<item><description>
<![CDATA[AUD 1.54220000 KRW 1354.37000000 NOK 10.24510000 ]]>
</description></item>`

	var input rss.Item

	assertions := assert.New(t)
	assertions.NoError(xml.Unmarshal([]byte(xmlResponse), &input))
	assertions.Len(input.Currencies.Rates, 3)
	assertions.Equal(input.Currencies.Rates[0], rss.Rate{Code: "AUD", Amount: "1.54220000"})
	assertions.Equal(input.Currencies.Rates[1], rss.Rate{Code: "KRW", Amount: "1354.37000000"})
	assertions.Equal(input.Currencies.Rates[2], rss.Rate{Code: "NOK", Amount: "10.24510000"})
}
