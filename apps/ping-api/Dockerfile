FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/ping-api ./

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /out/ping-api /app/ping-api
EXPOSE 8080
ENV PING_API_ADDR=:8080
CMD ["/app/ping-api"]
