FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata upx
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o pack-calculator .

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/pack-calculator /usr/local/bin/pack-calculator
COPY packs.json ./packs.json
COPY public ./public

ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/pack-calculator"]