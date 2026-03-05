FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o smart_home .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/smart_home .
EXPOSE 8080
CMD ["./smart_home"]