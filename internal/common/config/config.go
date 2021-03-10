package config

import (
	"flag"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Properties struct {
	Database DatabaseProperties
	Rss      RssProperties
}

type RssProperties struct {
	URL            string
	RequestTimeout time.Duration
}

type DatabaseProperties struct {
	ConnectionURL string
}

func LoadProperties() Properties {
	configPath := flag.String("config", "config/application-local.yaml", "")
	flag.Parse()

	viper.SetConfigFile(*configPath)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	var properties Properties

	err = viper.Unmarshal(&properties)
	if err != nil {
		log.Fatal(err)
	}

	return properties
}
