FROM golang:1.23-bullseye as base

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 65532 \
    small-user

WORKDIR /app

# Copy go mod and sum files first for caching
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Now copy the rest of the app
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main ./cmd/api

# Final image
FROM alpine

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group
COPY --from=base /main .

USER small-user:small-user

CMD ["./main"]
