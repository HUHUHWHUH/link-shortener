FROM golang:1.23-alpine AS builder

WORKDIR /shortener
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o shortenerbin ./cmd/shortener

FROM alpine:3.17
WORKDIR /shortener
COPY --from=builder /shortener/shortenerbin /shortener/shortenerbin
EXPOSE 8000

CMD ["/shortener/shortenerbin", "-storage", "memory"]