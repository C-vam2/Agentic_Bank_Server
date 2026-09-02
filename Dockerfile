# Build stage
FROM golang:1.27-alpine AS builder
WORKDIR /app
COPY . .

# Install migrate CLI
RUN go install -tags "postgres" \
    github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go build -o main main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env app.env
COPY db/migrations ./migrations
COPY start.sh .

# Copy migrate CLI
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]