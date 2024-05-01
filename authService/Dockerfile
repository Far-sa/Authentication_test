
FROM golang:1.19-alpine3.15 As builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Install MySQL and RabbitMQ client libraries
# RUN apk add --no-cache mariadb-client rabbitmq-c

COPY . .

RUN go build -o /app/auth-svc

# Runtime Steps
FROM alpine

COPY --from=builder /app/auth-svc .

EXPOSE 5001

CMD ["./auth-svc"]