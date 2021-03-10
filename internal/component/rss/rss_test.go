package rss_test

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/robertst0/ltc-currency-service/internal/common/config"
	"github.com/robertst0/ltc-currency-service/internal/component/currency"
	"github.com/robertst0/ltc-currency-service/internal/component/rss"
)

type RssTestSuite struct {
	assert *assert.Assertions
	suite.Suite
	DB         *sql.DB
	mock       sqlmock.Sqlmock
	service    *rss.Service
	mockServer *echo.Echo
}

func (suite *RssTestSuite) SetupSuite() {
	var err error

	suite.DB, suite.mock, err = sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      suite.DB,
		SkipInitializeWithVersion: true,
	})

	gormDB, err := gorm.Open(dialector)
	if err != nil {
		log.Fatal(err)
	}

	properties := config.RssProperties{RequestTimeout: 5 * time.Second, URL: "http://localhost:4567/rss"}
	suite.service = rss.NewService(currency.NewRepository(gormDB), properties)
	suite.mockServer = MockServer()

	go func() {
		suite.mockServer.Start(":4567")
	}()
}

func (suite *RssTestSuite) SetupTest() {
	suite.assert = assert.New(suite.T())
}

func (suite *RssTestSuite) TearDownSuite() {
	err := suite.mockServer.Shutdown(context.Background())
	if err != nil {
		log.Panic(err)
	}

	suite.DB.Close()
}

func TestRssTestSuite(t *testing.T) {
	suite.Run(t, new(RssTestSuite))
}

func (suite *RssTestSuite) TestService_FetchRssFeed() {
	suite.mockServer.GET("/rss", func(c echo.Context) error {
		return c.String(http.StatusOK, `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
	<channel>
    	<someIrrelevantTags>lorem ipsum</someIrrelevantTags>
		<item>
			<title>One</title>
			<description><![CDATA[AUD 1.54220000 ZAR 18.13530000 ]]></description>
			<pubDate>Tue, 02 Mar 2021 02:00:00 +0200</pubDate>
		</item>
		<item>
			<title>One</title>
			<description><![CDATA[ISK 152.50000000 TRY 8.96660000 ]]></description>
			<pubDate>Tue, 03 Mar 2021 02:00:00 +0200</pubDate>
		</item>
	</channel>
</rss>
`)
	})

	suite.mock.MatchExpectationsInOrder(true)

	suite.expectedInsert("AUD", "1.54220000", time.Date(2021, time.March, 2, 0, 0, 0, 0, time.UTC))
	suite.expectedInsert("ZAR", "18.13530000", time.Date(2021, time.March, 2, 0, 0, 0, 0, time.UTC))
	suite.expectedInsert("ISK", "152.50000000", time.Date(2021, time.March, 3, 0, 0, 0, 0, time.UTC))
	suite.expectedInsert("TRY", "8.96660000", time.Date(2021, time.March, 3, 0, 0, 0, 0, time.UTC))

	suite.service.RetrieveLatestData()

	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *RssTestSuite) TestService_FetchRssFeed_RetrieveError() {
	suite.mockServer.GET("/rss", func(c echo.Context) error {
		return c.String(http.StatusBadGateway, "Bad Gateway")
	})

	suite.service.RetrieveLatestData()

	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *RssTestSuite) expectedInsert(currency string, amount string, date time.Time) {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec("INSERT INTO `rate`").
		WithArgs(currency, amount, date).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()
}

func MockServer() *echo.Echo {
	return echo.New()
}
