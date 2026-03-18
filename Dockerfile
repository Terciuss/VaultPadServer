FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app ./cmd/server

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
