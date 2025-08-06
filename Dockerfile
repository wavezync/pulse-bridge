FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o pulse-bridge .

FROM alpine:latest AS runner

ENV GO_ENV=production

WORKDIR /app

COPY --from=builder /app/pulse-bridge .

COPY --from=builder /app/config.yml ./config.yml

RUN mkdir -p /etc/pulse-bridge && \
    cp ./config.yml /etc/pulse-bridge/config.yml && \
    chmod +x ./pulse-bridge

EXPOSE 8080

ENTRYPOINT ["./pulse-bridge"]