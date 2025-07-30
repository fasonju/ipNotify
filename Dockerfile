FROM golang:1.20-alpine AS builder

LABEL org.opencontainers.image.source="https://github.com/fasonju/ipNotify"
LABEL org.opencontainers.image.description="Public Ip Change notification app"
LABEL org.opencontainers.image.licenses=MIT

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ipNotify ./cmd/ipNotify

FROM scratch

COPY --from=builder /app/ipNotify /ipNotify

ENTRYPOINT ["/myapp"]
