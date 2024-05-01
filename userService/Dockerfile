FROM golang:1.19-alpine3.15 As builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Install MySQL and RabbitMQ client libraries
# RUN apk add --no-cache mariadb-client rabbitmq-c

COPY . .

RUN go build -o /app/user-svc

# Runtime Steps
FROM alpine


COPY --from=builder /app/user-svc .

EXPOSE 5000

CMD ["./user-svc"]
