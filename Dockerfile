FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o athena .

FROM alpine:latest AS runner

ENV GO_ENV=production

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/athena .
# Copy the config file from the builder stage
COPY --from=builder /app/config.yml ./config.yml
# Also create a backup in a standard location
RUN mkdir -p /etc/athena && \
    cp ./config.yml /etc/athena/config.yml && \
    chmod +x ./athena

EXPOSE 8080

ENTRYPOINT ["./athena"]