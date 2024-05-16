
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
COPY config.yaml ./
COPY database/migrations ./

EXPOSE 5001

# Set the entrypoint to the script
COPY scripts/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Custom initialization steps (database creation, etc.)
# COPY scripts/init.sh /usr/local/bin/init.sh
# RUN chmod +x /usr/local/bin/init.sh
# RUN /usr/local/bin/init.sh

ENTRYPOINT ["docker-entrypoint.sh"]

CMD ["./auth-svc"]