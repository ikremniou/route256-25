include make/lint.mk
include make/build.mk
include make/coverage.mk

install:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/air-verse/air@latest

lint: cart-lint loms-lint notifier-lint comments-lint

build: cart-build loms-build notifier-build comments-build

run-all:
	docker compose -f ./deploy/docker-compose.all.yaml up --build

run-dev:
	@{ cd cart && make watch; } & \
	{ cd loms && make watch; } & \
	{ docker compose -f ./deploy/docker-compose.dev.yaml up; } & \
	wait

run-dev-deps:
	docker compose -f ./deploy/docker-compose.dev.yaml up

coverage: cart-coverage

test-e2e:
	hurl --test ./cart/tests ./loms/tests

lint-cyclomatic:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

	@OUTPUT="$(shell gocyclo -over 10 -ignore '_test|pb.gw.go|pg.go' ./cart ./loms ./comments | tr '\n' '\1')"; \
	echo $$OUTPUT; \
	if [ -n "$$OUTPUT" ]; then \
		echo "\033[31mFailed to validation project functions cyclomatic complexity:\n\033[0m"; \
		echo "\033[31m$$OUTPUT\033[0m" | tr '\1' '\n'; \
		exit 1; \
	else \
		echo "\033[32mPassed cyclomatic complexity validation\033[0m"; \
	fi

lint-cognite:
	go install github.com/uudashr/gocognit/cmd/gocognit@latest

	@OUTPUT="$(shell gocognit -over 10 -ignore '_test|pb.gw.go|pg.go' ./cart ./loms ./comments | tr '\n' '\1')"; \
	if [ -n "$$OUTPUT" ]; then \
		echo "\033[31mFailed to validation project functions cognitive complexity:\n\033[0m"; \
		echo "\033[31m$$OUTPUT\033[0m" | tr '\1' '\n'; \
		exit 1; \
	else \
		echo "\033[32mPassed cognite complexity validation\033[0m"; \
	fi
