.DEFAULT_GOAL := build

.PHONY: prerequisite_check
prerequisite_check:
	@echo Checking prerequisites.
	@make -v | head -n 1
	@docker -v
	@docker-compose -v
	@echo Prerequisote check passed

.PHONY: api_app
api_app: prerequisite_check
	docker build --rm -t ltc/ltc-homework-api_app -f build/docker/api/Dockerfile .

.PHONY: rss_app
rss_app: prerequisite_check
	docker build --rm -t ltc/ltc-homework-rss_app -f build/docker/rss/Dockerfile .

.PHONY: build
build: api_app rss_app

.PHONY: start
start:
	docker-compose up -d

.PHONY: fetch_rss
fetch_rss:
	docker run --network=ltc-currency-service_ltc ltc/ltc-homework-rss_app

.PHONY: stop
stop:
	docker-compose stop
