.PHONY: build

cart-build:
	cd cart && GOOS=linux GOARCH=amd64 make build

loms-build:
	cd loms     && GOOS=linux GOARCH=amd64 make build

notifier-build:
	cd notifier && GOOS=linux GOARCH=amd64 make build

comments-build:
	cd comments && GOOS=linux GOARCH=amd64 make build
