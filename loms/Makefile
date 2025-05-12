BINDIR=${CURDIR}/bin
PACKAGE=route256/loms
LOCAL_CONFIG=${CURDIR}/configs/values_local.yaml
BUILD_COMMAND=go build -o ${BINDIR}/loms_service cmd/loms_service/server.go
VENDOR_BIN=${CURDIR}/vendor/bin

DB_HOST ?= localhost
LOMS_DB_USER ?= loms-user
LOMS_DB_PASS ?= loms-password
LOMS_DB_PORT ?= 5432
LOMS_DB_NAME ?= loms_db
LOMS_DB_CONN_STRING = postgres://$(LOMS_DB_USER):$(LOMS_DB_PASS)@$(DB_HOST):$(LOMS_DB_PORT)/$(LOMS_DB_NAME)?sslmode=disable
MIGRATION_DIR ?= ./migrations

bindir:
	mkdir -p ${BINDIR}

build: bindir
	${BUILD_COMMAND}

run: build
	export CONFIG_FILE=${LOCAL_CONFIG} && \
	${BINDIR}/loms_service

watch:
	export CONFIG_FILE=${LOCAL_CONFIG} && air --build.cmd "${BUILD_COMMAND}" \
		--build.bin ${BINDIR}/loms_service --build.send_interrupt true --build.kill_delay 500

generate:
	buf generate
	sqlc generate

serve-swagger:
	docker run -p 8086:8080 --name loms-swagger --rm \
		-e SWAGGER_JSON=/app/apidocs.swagger.json \
		-v ./api/openapiv2/apidocs.swagger.json:/app/apidocs.swagger.json \
		swaggerapi/swagger-ui

migrate-create:
	mkdir -p $(MIGRATION_DIR)
	goose -dir ${MIGRATION_DIR} create $(word 2, $(MAKECMDGOALS)) sql
	
migrate-up:
	goose -dir ${MIGRATION_DIR} postgres "${LOMS_DB_CONN_STRING}" up

migrate-down:
	goose -dir ${MIGRATION_DIR} postgres "${LOMS_DB_CONN_STRING}" down

# runs migrations on CI
run-migrations: migrate-up