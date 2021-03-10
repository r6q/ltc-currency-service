package api_test

import (
	"database/sql"
	"log"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/robertst0/ltc-currency-service/internal/component/api"
	"github.com/robertst0/ltc-currency-service/internal/component/currency"
)

var now = time.Date(2021, time.March, 7, 0, 0, 0, 0, time.UTC)

type APITestSuite struct {
	assert *assert.Assertions
	suite.Suite
	DB   *sql.DB
	mock sqlmock.Sqlmock
	app  *echo.Echo
}

func (suite *APITestSuite) SetupSuite() {
	suite.app = echo.New()

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

	handler := api.NewHandler(api.NewService(currency.NewRepository(gormDB)))
	handler.Route(suite.app.Group("/api/v1/rates"))
}

func (suite *APITestSuite) SetupTest() {
	suite.assert = assert.New(suite.T())
}

func (suite *APITestSuite) TearDownSuite() {
	suite.DB.Close()
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (suite *APITestSuite) TestService_RetrieveAllRates() {
	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rate`")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "currency_code", "amount", "date"}).
			AddRow(1, "GBP", "0.86433000", now.Add(-24*time.Hour)).
			AddRow(2, "USD", "1.20290000", now.Add(-24*time.Hour)).
			AddRow(3, "GBP", "0.86351000", now).
			AddRow(4, "USD", "1.20330000", now))

	apitest.New().
		Handler(suite.app).
		Get("/api/v1/rates").
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(jsonpath.Len("$.data", 4)).
		Assert(jsonpath.Root("$.data[0]").
			Equal("currency_code", "GBP").
			Equal("amount", "0.86433000").
			Equal("date", "2021-03-06T00:00:00Z").
			End()).
		Assert(jsonpath.Root("$.data[1]").
			Equal("currency_code", "USD").
			Equal("amount", "1.20290000").
			Equal("date", "2021-03-06T00:00:00Z").
			End()).
		Assert(jsonpath.Root("$.data[2]").
			Equal("currency_code", "GBP").
			Equal("amount", "0.86351000").
			Equal("date", "2021-03-07T00:00:00Z").
			End()).
		Assert(jsonpath.Root("$.data[3]").
			Equal("currency_code", "USD").
			Equal("amount", "1.20330000").
			Equal("date", "2021-03-07T00:00:00Z").
			End()).
		End()
}

func (suite *APITestSuite) TestService_RetrieveLatestRates() {
	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rate` WHERE date = (SELECT MAX(date) FROM `rate`)")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "currency_code", "amount", "date"}).
			AddRow(3, "GBP", "0.86351000", now).
			AddRow(4, "USD", "1.20330000", now))

	apitest.New().
		Handler(suite.app).
		Get("/api/v1/rates").
		Query("latest", "true").
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(jsonpath.Len("$.data", 2)).
		Assert(jsonpath.Root("$.data[0]").
			Equal("currency_code", "GBP").
			Equal("amount", "0.86351000").
			Equal("date", "2021-03-07T00:00:00Z").
			End()).
		Assert(jsonpath.Root("$.data[1]").
			Equal("currency_code", "USD").
			Equal("amount", "1.20330000").
			Equal("date", "2021-03-07T00:00:00Z").
			End()).
		End()
}

func (suite *APITestSuite) TestService_RetrieveLatestRates_WithInvalidParam() {
	apitest.New().
		Handler(suite.app).
		Get("/api/v1/rates").
		Query("latest", "something").
		Expect(suite.T()).
		Status(http.StatusBadRequest).
		Assert(jsonpath.Equal("$.error", "Param value 'something' is invalid")).
		End()
}

func (suite *APITestSuite) TestService_RetrieveHistoricalRates() {
	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rate` WHERE currency_code = ?")).
		WithArgs("USD").
		WillReturnRows(sqlmock.NewRows([]string{"id", "currency_code", "amount", "date"}).
			AddRow(2, "USD", "1.20290000", now.Add(-24*time.Hour)).
			AddRow(4, "USD", "1.20330000", now))

	apitest.New().
		Handler(suite.app).
		Get("/api/v1/rates").
		Query("currency", "USD").
		Expect(suite.T()).
		Status(http.StatusOK).
		Assert(jsonpath.Len("$.data", 2)).
		Assert(jsonpath.Root("$.data[0]").
			Equal("currency_code", "USD").
			Equal("amount", "1.20290000").
			Equal("date", "2021-03-06T00:00:00Z").
			End()).
		Assert(jsonpath.Root("$.data[1]").
			Equal("currency_code", "USD").
			Equal("amount", "1.20330000").
			Equal("date", "2021-03-07T00:00:00Z").
			End()).
		End()
}
