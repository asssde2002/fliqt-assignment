FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o db_migrate ./cmd/migrate/main.go

FROM alpine
WORKDIR /app
RUN apk add --no-cache ca-certificates bash tzdata
COPY --from=builder /app/main /app/db_migrate /app/scripts /app/
