DOCKER = docker
GO = go

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
