# Commands MUST be indented using a tab; using spaces will not work.
# Using a @ will suppress "make" from echoing out the command when it's ran.
# Using .PHONE: <command> tells "make" that the name is something that should be executed, and not a file.

# Include variables from the .envrc file
include .envrc

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
# DEVELOPMENT
# ==================================================================================== #

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DB_CONN}

## db/mig/new name=$1: create a new database migration
.PHONY: db/mig/new
db/mig/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

# note: the $DB_CONN string is being sourced from the .envrc file.
## db/mig/up: apply all up database migrations
.PHONY: db/mig/up
db/mig/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${DB_CONN} up

## dev/web: run the cmd/web application using 'air' for live reload
.PHONY: dev/web
dev/web:
	@air

## run/docker/web: run the dockerized cmd/web application
.PHONY: run/docker/web
run/docker/web: build/web
	@echo 'Starting docker container for cmd/web...'
	podman run -p 4000:4000 --rm localhost/${APP_DOCKER_NAME}:latest

## run/web: run the cmd/web application
.PHONY: run/web
run/web:
	@go run ./cmd/web

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit:
	@echo 'Checking module dependencies'
	@#go mod tidy -diff | only applicable in Go 1.23 and later
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## tidy: format all .go files and tidy module dependencies
.PHONY: tidy
tidy:
	@echo 'Formatting .go files...'
	go fmt ./...
	@echo 'Tidying module dependencies...'
	go mod tidy
	@echo 'Verifying and vendoring module dependencies...'
	go mod verify
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/docker/web: build the cmd/web application in a podman/docker container
.PHONY: build/docker/web
build/docker/web: build/web
	@echo 'Building docker container for cmd/web...'
	podman build -t ${APP_DOCKER_NAME} .

## build/web: build the cmd/web application
.PHONY: build/web
build/web:
	@echo 'Building cmd/web...'
	CGO_ENABLED=0 go build -ldflags='-s' -o=./bin/host_web ./cmd/web
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64_web ./cmd/web
