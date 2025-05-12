BINDIR=${CURDIR}/bin
PACKAGE=route256/cart
LOCAL_CONFIG=${CURDIR}/configs/values_local.yaml
BUILD_COMMAND=go build -o ${BINDIR}/cart_service cmd/cart_service/server.go

.PHONY: run test coverage

bindir:
	mkdir -p ${BINDIR}

build: bindir
	${BUILD_COMMAND}

run:
	export CONFIG_FILE=${LOCAL_CONFIG} && go run ./cmd/cart_service

watch:
	export CONFIG_FILE=${LOCAL_CONFIG} && air --build.cmd "${BUILD_COMMAND}" \
		--build.bin ${BINDIR}/cart_service --build.send_interrupt true --build.kill_delay 500

test:
	go test ./...

test-race:
	go test ./... -race

generate:
	buf generate

coverage-short:
	@go test ./... -coverprofile=.coverage > /dev/null
	@go tool cover -func=.coverage | tail -n 1 | awk '{print $$3}' | xargs -I {} echo "Total project 'cart' coverage is: {}"
	@rm -f .coverage
