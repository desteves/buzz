# syntax=docker/dockerfile:1
FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o buzz

FROM alpine AS release
WORKDIR /app
COPY --from=builder /app/buzz buzz
COPY --from=builder /app/static static
EXPOSE 8000
ENV SERVER_ADDR=0.0.0.0:8000
ENTRYPOINT ["./buzz"]