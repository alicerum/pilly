DOCKER = docker
GO = go
CONFIG = ./config.yaml

test-db:
	$(DOCKER) run \
		--name postgres \
		-p 5432:5432 \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=pilly \
		-d postgres

build:
	$(GO) build ./cmd/pilly

.PHONY: migrations
migrate:
	$(GO) run ./cmd/migrations/*.go -config $(CONFIG)

.PHONY: migrations
migrate-down:
	$(GO) run ./cmd/migrations/*.go -config $(CONFIG) down

