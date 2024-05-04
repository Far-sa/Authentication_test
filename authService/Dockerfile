
FROM golang:1.19-alpine3.15 As builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /app/auth-svc

# Runtime Steps
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-svc .

# Explicitly copy the configuration file into the container
# COPY auth-svc/config ./config

EXPOSE 5001

ENV RABBITMQ_HOST rabbitmq
ENV RABBITMQ_USER guest
ENV RABBITMQ_PASSWORD guest
ENV RABBITMQ_PORT 5672

CMD ["./auth-svc"]