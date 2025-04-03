FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o pulsebridge .

FROM alpine:latest AS runner

ENV GO_ENV=production

WORKDIR /app

COPY --from=builder /app/pulsebridge .

COPY --from=builder /app/config.yml ./config.yml

RUN mkdir -p /etc/pulsebridge && \
    cp ./config.yml /etc/pulsebridge/config.yml && \
    chmod +x ./pulsebridge

EXPOSE 8080

ENTRYPOINT ["./pulsebridge"]