FROM golang:1.23-alpine AS builder

# Install make
RUN apk add --no-cache make

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN make build


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/loms_service ./
COPY configs/values_docker.yaml ./configs/values.yaml

ENV CONFIG_FILE=/app/configs/values.yaml
EXPOSE 8083
EXPOSE 8084

CMD ["./loms_service"]