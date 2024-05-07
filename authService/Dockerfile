
FROM golang:1.19-alpine3.15 As builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /app/auth-svc

# Runtime Steps
FROM alpine:latest


COPY --from=builder /app/auth-svc .

# Explicitly copy the configuration file into the container
COPY config.yaml ./
COPY database/migrations ./

EXPOSE 5001

CMD ["./auth-svc"]