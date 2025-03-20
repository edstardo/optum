
FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/bin/service ./cmd/market-data-service

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /app/bin/service /app/service

USER nonroot:nonroot

ENTRYPOINT ["/app/service"]
