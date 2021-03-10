# LTC Currency Service

[![effective go](https://img.shields.io/badge/code_style-%E2%9D%A4-FF4081.svg)](https://golang.org/doc/effective_go.html)

### Development

To run tests with coverage use the following command:
```
go test ./... -coverpkg=./... -coverprofile=coverage.out ./... -count=1
```
_To see coverage results in browser use the following command:_
```
go tool cover -html=coverage.out
```

### Build

To build applications simply invoke `make` command on the root directory of the project

Makefile will create docker image for each application

### Run application

To start the api service run the following command:
```
make start
```
It will start `MySql` database and `API` service afterwards. Alternatively you can use `docker-compose up`

In order for api service to give out currency rates, you need to fetch the data from external service.
To do so run the following command:
```
make fetch_rss
```

To stop the api service run the following command:
```
make stop
```

### Endpoints

The API service provides 3 endpoints for data retrieval:

- http://localhost:8080/api/v1/rates (retrieves all rates)
- http://localhost:8080/api/v1/rates?latest=true (retrieves only latest rates)
- http://localhost:8080/api/v1/rates?currency=USD (retrieves all rates for specified currency)

Alternatively if you are using Intellij IDEA or GoLand there is a `.http` file under `integration` 
folder where all commands are conveniently laid out

### Linter

Application uses golangci linter which consists of 57 linters

To run linter, execute the following command:

* If installed on the machine:
```
golangci-lint run
``` 

* If downloaded docker image:
```
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.37.0 golangci-lint run -v
``` 


### Some design decisions

I decided to have two standalone services where one is dedicated API service, and the other
one is short-lived data retrieval service.

The idea behind having dedicated RSS retrieval service is to imitate serverless approach where
it can be invoked by some scheduled cron job or event of sorts, do its job by retrieving data
and become dormant again 

TODO:
- [ ] Add proper `logger` and combine it with `echo` logging middleware to have unified logger everywhere
- [ ] Add a couple of more tests for some very unlikely test scenarios
- [ ] Add an external database changeset tool instead of relying on auto-migration with GORM 
