FROM golang:1.24

LABEL org.opencontainers.image.source="https://github.com/fasonju/ipNotify"
LABEL org.opencontainers.image.description="Public Ip Change notification app"
LABEL org.opencontainers.image.licenses=MIT

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /app/ipNotify ./cmd/ipNotify

CMD ["./ipNotify"]
