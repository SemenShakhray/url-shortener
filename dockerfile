FROM golang:1.23.4-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash git make gettext musl-dev

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# build
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
RUN CGO_ENABLED=0 go build -o /app/shortener ./cmd/url-shortener/main.go

FROM alpine AS runner

WORKDIR /app

COPY --from=builder /app/shortener /app/shortener
COPY .env /app/.env
COPY /migrations /app/migrations

CMD ["/app/shortener"]