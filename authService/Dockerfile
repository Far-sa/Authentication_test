
FROM golang:1.22-alpine As builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /app/auth-svc

# Runtime Steps
FROM alpine:latest

# WORKDIR /app

COPY --from=builder /app/auth-svc .

# Explicitly copy the configuration file into the container
COPY config.yaml ./
COPY database/migrations ./

# Custom initialization steps (database creation, etc.)
# COPY scripts/init-db.sh /docker-entrypoint-initdb.d/init-db.sh
# RUN chmod +x /docker-entrypoint-initdb.d/init-db.sh

EXPOSE 5001

# Install required packages (including make) and run make command
# RUN apk --no-cache add make \
#     && make init-db

# ENTRYPOINT ["/docker-entrypoint-initdb.d/init-db.sh"]

CMD ["./auth-svc"]