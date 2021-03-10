package rss

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const DateTimeFormat = "Mon, 02 Jan 2006 15:04:05 -0700"

type Response struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Currencies Currency `xml:"description"`
	Date       DateTime `xml:"pubDate"`
}

type DateTime struct {
	time.Time
}

func (d *DateTime) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var value string

	err := decoder.DecodeElement(&value, &start)
	if err != nil {
		return fmt.Errorf("decoding failed: %w", err)
	}

	dateTime, err := time.Parse(DateTimeFormat, value)
	if err != nil {
		return fmt.Errorf("decoding failed: %w", err)
	}

	*d = DateTime{dateTime.UTC()}

	return nil
}

type Currency struct {
	Rates []Rate
}

type Rate struct {
	Code   string
	Amount string
}

func (c *Currency) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var value string

	err := decoder.DecodeElement(&value, &start)
	if err != nil {
		return fmt.Errorf("decoding failed: %w", err)
	}

	pairs := regexp.MustCompile(`[A-Z]{3} \d*\.\d+`).FindAllString(value, -1)

	rates := make([]Rate, len(pairs))

	for i, pair := range pairs {
		rate := strings.Split(pair, " ")
		rates[i] = Rate{
			Code:   rate[0],
			Amount: rate[1],
		}
	}

	*c = Currency{rates}

	return nil
}
