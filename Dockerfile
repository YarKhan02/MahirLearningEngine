FROM golang:1.26.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o server \
    ./cmd/server

FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /app/server .
# The server runs migrations at boot (MIGRATIONS_PATH must point here).
COPY --from=builder /app/migrations /migrations

USER nonroot:nonroot

ENTRYPOINT ["/server"]
