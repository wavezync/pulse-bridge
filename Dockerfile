FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o pulsebridge .

FROM alpine:latest AS runner

ENV GO_ENV=production

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/pulsebridge .
# Copy the config file from the builder stage
COPY --from=builder /app/config.yml ./config.yml
# Also create a backup in a standard location
RUN mkdir -p /etc/pulsebridge && \
    cp ./config.yml /etc/pulsebridge/config.yml && \
    chmod +x ./pulsebridge

EXPOSE 8080

ENTRYPOINT ["./pulsebridge"]