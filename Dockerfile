FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o monsite

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/monsite .


COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["./monsite"]

