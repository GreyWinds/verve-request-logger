
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o verve main.go


FROM alpine:latest
WORKDIR /root/
RUN apk --no-cache add redis
COPY --from=builder /app/verve .
RUN touch unique_requests.log
EXPOSE 8080

CMD ["./verve"]
