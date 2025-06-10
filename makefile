# Commands MUST be indented using a tab; using spaces will not work.
# Using a @ will suppress "make" from echoing out the command when it's ran.
# Using .PHONE: <command> tells "make" that the name is something that should be executed, and not a file.

# Include variables from the .env file
include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# Confirm that the user wants to run the make command.
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DATABASE
# ==================================================================================== #

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	@psql ${DB_CONN}

## db/mig/create name=$1: create a new database migration
.PHONY: db/mig/create
db/mig/create:
	@echo 'Creating migration files for ${name}...'
	@goose -dir=./migrations postgres ${DB_CONN} create ${name} sql

## db/mig/down: apply all up database migrations
.PHONY: db/mig/down
db/mig/down: confirm
	@echo 'Rolling back last migration...'
	@goose -dir=./migrations postgres ${DB_CONN} down

## db/mig/status: see migration status for current database
.PHONY: db/mig/status
db/mig/status:
	@echo 'Getting migration status for database...'
	@goose -dir=./migrations postgres ${DB_CONN} status

## db/mig/up: apply all up database migrations
.PHONY: db/mig/up
db/mig/up: confirm
	@echo 'Running up migrations...'
	@goose -dir=./migrations postgres ${DB_CONN} up

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit:
	@echo 'Checking module dependencies'
	@go mod tidy -diff
	@go mod verify
	@echo 'Vetting code...'
	@go vet ./...
	@staticcheck ./...
	@echo 'Running tests...'
	@go test -race -vet=off ./...

## tidy: format all .go files and tidy module dependencies
.PHONY: tidy
tidy:
	@echo 'Tidying module dependencies...'
	@go mod tidy
	@echo 'Verifying and vendoring module dependencies...'
	@go mod verify
	@go mod vendor
	@echo 'Formatting .go files...'
	@go fmt ./...

## up: update all dependencies in go.mod and runs make tidy afterwards
.PHONY: up
up:
	@echo 'Updating dependencies...'
	@go get -u ./...
	$(MAKE) tidy

## test: run all the tests
.PHONY: test
test: test/internal test/web

## test/internal: run test in ./internal
.PHONY: test/internal
test/internal:
	@echo 'Running tests in ./internal...'
	@go test -race -vet=off ./internal/...

## test/web: run test in ./cmd/web
.PHONY: test/web
test/web:
	@echo 'Running tests in ./cmd/web...'
	@go test -race -vet=off ./cmd/web/...

# ==================================================================================== #
# DOCKER
# ==================================================================================== #

## docker/build/web: build the cmd/web application in a docker container
.PHONY: docker/build/web
docker/build/web:
	@echo 'Building docker container for cmd/web...'
	sudo docker build --no-cache -t ${APP_DOCKER_NAME} . && sudo docker image prune -f

## docker/run/web: run the dockerized cmd/web application
.PHONY: docker/run/web
docker/run/web: docker/build/web
	@echo 'Starting docker container for cmd/web...'
	sudo docker run -p 4000:4000 --network="host" --env-file=".env" ${APP_DOCKER_NAME}

# ==================================================================================== #
# WEB
# ==================================================================================== #

## web/build: build the cmd/web application
.PHONY: web/build
web/build:
	@echo 'Building cmd/web...'
	CGO_ENABLED=0 go build -ldflags='-s' -o=./bin/host_web ./cmd/web
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64_web ./cmd/web

## web/dev: run the cmd/web application using 'air' for live reload
.PHONY: web/dev
web/dev:
	@air

## web/run: run the cmd/web application (fallback if 'air' isn't installed)
.PHONY: web/run
web/run:
	@go run ./cmd/web \
		-dsn=${DB_CONN}
